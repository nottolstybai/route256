package app

import (
	"context"
	"github.com/IBM/sarama"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"route256.ozon.ru/project/loms/internal/app/outbox"
	"route256.ozon.ru/project/loms/internal/app/server"
	"route256.ozon.ru/project/loms/internal/config"
	"route256.ozon.ru/project/loms/internal/infra/kafka/producer"
	"route256.ozon.ru/project/loms/internal/mw"
	"route256.ozon.ru/project/loms/internal/repository/db"
	producer_handler "route256.ozon.ru/project/loms/internal/repository/kafka/producer"
	service "route256.ozon.ru/project/loms/internal/service/loms"
	outbox_service "route256.ozon.ru/project/loms/internal/service/outbox"
	desc "route256.ozon.ru/project/loms/pkg/api/loms/v1"
	"route256.ozon.ru/project/loms/pkg/logger"
	"time"
)

type ConnectionCloser interface {
	CloseConnections() error
}

type App struct {
	config        config.Config
	gateway       *http.Server
	gatewayMux    *runtime.ServeMux
	server        *grpc.Server
	storageCloser ConnectionCloser

	outbox         *outbox.OutboxSender
	producerCloser ConnectionCloser
}

func NewApp(config config.Config) *App {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(),
		grpcprom.WithServerCounterOptions(),
	)
	prometheus.DefaultRegisterer.MustRegister(srvMetrics)

	// create grpc server object
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.Trace,
			srvMetrics.UnaryServerInterceptor(),
			mw.Panic,
			mw.Logger,
			mw.Validate,
		),
	)
	reflection.Register(grpcServer)

	// define dependencies
	ctx := context.Background()
	masterPool, replicaPool := SetupDBConnection(ctx, config.DBMasterConnString, config.DBReplicaConnString)

	// create sync producer for kafka
	syncProducer := SetupProducer(config)
	producerHandler := producer_handler.NewHandler(syncProducer, config.KafkaConfig.Producer.Topic)

	// create storage for services
	lomsStorage := db.NewLomsStorageFromJson(masterPool, replicaPool, producerHandler)

	// create outbox service
	outboxService := outbox_service.NewOutboxService(lomsStorage)
	outboxSender := outbox.NewOutboxSender(outboxService)

	// create loms service and its controller
	lomsService := service.NewLOMSService(lomsStorage)
	controller := server.NewServer(lomsService)

	// register server
	desc.RegisterLomsServer(grpcServer, controller)

	// create mux for http gateway
	mux := http.NewServeMux()
	gwMux := runtime.NewServeMux()

	mux.Handle("/", gwMux)
	mux.Handle("/metrics", promhttp.Handler())

	// create http server
	gwServer := &http.Server{
		Addr:    config.GatewayAddr,
		Handler: mw.WithHTTPLoggingMiddleware(mux),
	}

	return &App{
		config:         config,
		gateway:        gwServer,
		gatewayMux:     gwMux,
		server:         grpcServer,
		storageCloser:  lomsStorage,
		outbox:         outboxSender,
		producerCloser: producerHandler,
	}
}

func (a *App) grpcServe() {
	// create listener for grpc server
	conn, err := net.Listen("tcp", a.config.ServeAddr)
	if err != nil {
		logger.Fatal("failed create listener", zap.Error(err))
	}
	defer conn.Close()

	// serve grpc server
	logger.Info("grpc server listening", zap.Any("address", conn.Addr()))
	if err = a.server.Serve(conn); err != nil {
		logger.Fatal("failed to serve grpc server", zap.Error(err))
	}
}

func (a *App) gatewayServe() {
	// setup conn to grpc server
	conn, err := grpc.Dial(a.config.ServeAddr, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed create listener", zap.Error(err))
	}
	defer conn.Close()

	//register all handlers of gateway
	if err := desc.RegisterLomsHandler(context.Background(), a.gatewayMux, conn); err != nil {
		logger.Fatal("failed to register gateway", zap.Error(err))
	}

	// serve gateway
	logger.Info("serving gRPC-Gateway", zap.String("address", a.config.GatewayAddr))
	if err := a.gateway.ListenAndServe(); err != nil {
		logger.Fatal("failed to serve gateway", zap.Error(err))
	}
}

func (a *App) Run(ctx context.Context) {
	go a.grpcServe()
	go a.gatewayServe()
	go a.outbox.RunDispatcher(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	a.server.GracefulStop()

	if err := a.storageCloser.CloseConnections(); err != nil {
		return err
	}
	if err := a.producerCloser.CloseConnections(); err != nil {
		return err
	}
	return a.gateway.Shutdown(ctx)
}

// SetupDBConnection sets the connection to master and replica dbs.
// If we were not able to connect to master throw panic
// If we were not able to connect to replica return only connections to master
func SetupDBConnection(ctx context.Context, masterDSN, replicaDSN string) (*pgxpool.Pool, *pgxpool.Pool) {
	masterPool, err := pgxpool.New(ctx, masterDSN)
	if err != nil {
		panic(err)
	}
	if err := masterPool.Ping(ctx); err != nil {
		panic(err)
	}

	replicaPool, err := pgxpool.New(ctx, replicaDSN)
	if err != nil {
		logger.Warn("couldn't connect to replica", zap.Error(err))
		return masterPool, masterPool
	}
	if err = replicaPool.Ping(ctx); err != nil {
		logger.Warn("couldn't connect to replica", zap.Error(err))
		return masterPool, masterPool
	}

	return masterPool, replicaPool
}

// SetupProducer creates a sync producer for kafka
func SetupProducer(cfg config.Config) sarama.SyncProducer {
	prod, err := producer.NewSyncProducer(cfg.KafkaConfig.Kafka,
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(10),
		producer.WithRetryBackoff(5*time.Millisecond),
		producer.WithProducerPartitioner(sarama.NewRoundRobinPartitioner),
	)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return prod
}

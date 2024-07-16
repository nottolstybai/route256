package config

import (
	"os"
	"route256.ozon.ru/project/loms/internal/infra/kafka"
	"route256.ozon.ru/project/loms/internal/repository/kafka/producer"
)

const (
	defaultGRPCHost        = ":50051"
	defaultHTTPHost        = ":3000"
	defaultDBMasterHost    = "postgres://user:password@localhost:5432/postgres"
	defaultDBReplicaHost   = "postgres://user:password@localhost:5433/postgres"
	defaultBootstrapServer = "localhost:9092"
	defaultKafkaTopic      = "loms.order-events"
)

type KafkaConfig struct {
	Kafka    kafka.Config
	Producer producer.Config
}

type Config struct {
	ServeAddr           string
	GatewayAddr         string
	DBMasterConnString  string
	DBReplicaConnString string
	KafkaConfig         KafkaConfig
}

func NewConfig() Config {
	return Config{
		ServeAddr:           getEnvHelper("LOMS_HOST_ADDR", defaultGRPCHost),
		GatewayAddr:         getEnvHelper("LOMS_GATEWAY_ADDR", defaultHTTPHost),
		DBMasterConnString:  getEnvHelper("POSTGRES_MASTER_CONN_STR", defaultDBMasterHost),
		DBReplicaConnString: getEnvHelper("POSTGRES_REPLICA_CONN_STR", defaultDBReplicaHost),
		KafkaConfig: KafkaConfig{
			Kafka: kafka.Config{
				Brokers: []string{
					getEnvHelper("BOOTSTRAP_SERVER", defaultBootstrapServer),
				},
			},
			Producer: producer.Config{Topic: getEnvHelper("KAFKA_TOPIC", defaultKafkaTopic)},
		},
	}
}

func getEnvHelper(envVar, defaultVal string) string {
	val := os.Getenv(envVar)
	if val == "" {
		val = defaultVal
	}
	return val
}

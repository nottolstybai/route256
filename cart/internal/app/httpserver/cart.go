package httpserver

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"route256.ozon.ru/project/cart/internal/entity"
	"route256.ozon.ru/project/cart/internal/mw"
	"route256.ozon.ru/project/cart/pkg/logger"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Service interface {
	AddItem(ctx context.Context, userID, skuID int64, count uint16) error
	DeleteItem(ctx context.Context, userID, skuID int64) error
	DeleteItemsByUserID(ctx context.Context, userID int64) error
	GetItemsByUserID(ctx context.Context, userID int64) (*entity.ListItems, error)
	Checkout(ctx context.Context, userID int64) (int64, error)
}

type CartServer struct {
	server *http.Server
	cart   Service
}

func NewServer(cart Service, address string) *CartServer {
	return &CartServer{cart: cart, server: &http.Server{Addr: address}}
}

func (s *CartServer) Serve() error {
	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", mw.HandleWithSpan(mw.HandleMetrics(mw.Log(s.AddItem))))
	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", mw.HandleWithSpan(mw.HandleMetrics(mw.Log(s.DeleteItem))))
	http.HandleFunc("DELETE /user/{user_id}/cart", mw.HandleWithSpan(mw.HandleMetrics(mw.Log(s.DeleteAllItems))))
	http.HandleFunc("GET /user/{user_id}/cart", mw.HandleWithSpan(mw.HandleMetrics(mw.Log(s.ListItems))))
	http.HandleFunc("POST /cart/checkout", mw.HandleWithSpan(mw.HandleMetrics(mw.Log(s.Checkout))))
	http.Handle("/metrics", promhttp.Handler())

	logger.Info("Starting server")

	return s.server.ListenAndServe()
}

func (s *CartServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

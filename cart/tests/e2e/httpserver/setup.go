package httpserver

import (
	"net/http"
	"route256.ozon.ru/project/cart/internal/app/httpserver"
	"route256.ozon.ru/project/cart/internal/client/product_service"
	"route256.ozon.ru/project/cart/internal/config"
	"route256.ozon.ru/project/cart/internal/repository/memory"
	"route256.ozon.ru/project/cart/internal/service"
)

func MakeServer() *httpserver.CartServer {
	cfg := config.NewConfig()

	repo := memory.NewMemoryRepo()
	productService := product_service.NewProductServiceClient(cfg.ProductServiceHost, cfg.ProductServiceToken)
	cartService := service.NewCartService(productService, repo)

	server := httpserver.NewServer(cartService, cfg.ServeAddr)
	return server
}

func MakeClient() *http.Client {
	return &http.Client{}
}

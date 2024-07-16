package config

import "os"

const (
	defaultAddress            = ":8082"
	defaultToken              = "testtoken"
	defaultProductServiceHost = "http://route256.pavl.uk:8080"
	defaultLomsServiceHost    = ":50051"
)

type Config struct {
	ServeAddr           string
	ProductServiceToken string
	ProductServiceHost  string
	LomsServiceHost     string
}

func NewConfig() *Config {
	return &Config{
		ServeAddr:           getEnvHelper("CART_HOST_ADDR", defaultAddress),
		ProductServiceToken: getEnvHelper("PRODUCT_SERVICE_TOKEN", defaultToken),
		ProductServiceHost:  getEnvHelper("PRODUCT_SERVICE_HOST", defaultProductServiceHost),
		LomsServiceHost:     getEnvHelper("LOMS_SERVICE_HOST", defaultLomsServiceHost),
	}
}

func getEnvHelper(envVar, defaultVal string) string {
	val := os.Getenv(envVar)
	if val == "" {
		val = defaultVal
	}
	return val
}

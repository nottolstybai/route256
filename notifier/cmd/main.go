package main

import (
	"route256.ozon.ru/project/notifier/internal/app"
	"route256.ozon.ru/project/notifier/pkg/logger"
)

func main() {
	logger.Init()
	defer logger.Sync()

	config := app.NewConfig()
	kafkaApp := app.NewApp(config)
	kafkaApp.Run()
}

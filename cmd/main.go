package main

import (
	"NotRedis/internal/app"
	"NotRedis/internal/compute/parser"
	"NotRedis/internal/storage"
	"NotRedis/internal/storage/engine"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	myMap := engine.NewMyMap(100, logger)
	storage, _ := storage.NewStorage(myMap, logger)
	parser := parser.NewParser(logger)

	app := app.NewApp(parser, storage, logger)
	app.Run()
}

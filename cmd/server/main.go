package main

import (
	"NotRedis/internal/compute/parser"
	"NotRedis/internal/config"
	"NotRedis/internal/database/engine"
	"NotRedis/internal/database/storage"
	"NotRedis/internal/network"
	"fmt"
	"os"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	logger, err := config.SetupLogger(cfg)
	if err != nil {
		fmt.Println("Error setting up logger:", err)
		os.Exit(1)
	}
	defer logger.Sync()

	myMap := engine.NewMyMap(100, logger)

	store, err := storage.NewStorage(myMap, logger)
	if err != nil {
		logger.Fatal("failed to create storage")
	}

	p := parser.NewParser(logger)

	serverCfg := config.ToServerConfig(cfg)
	srv := network.NewServer(serverCfg, store, p, logger)

	if err := srv.Start(); err != nil {
		logger.Fatal("server failed")
	}
}

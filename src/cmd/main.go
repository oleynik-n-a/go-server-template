package main

import (
	"time"

	"github.com/go-server-template/internal/repository"
	"github.com/go-server-template/internal/service"
	v1 "github.com/go-server-template/internal/transport/rest/v1"
	"github.com/go-server-template/pkg/database"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	postgresCfg, err := repository.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading postgres config", zap.Error(err))
		return
	}

	db, err := repository.ConnectPostgres(postgresCfg)
	if err != nil {
		logger.Fatal("Error connecting Database", zap.Error(err))
		return
	}
	logger.Info("Connected to PostgreSQL!", zap.String("port", postgresCfg.DBPort))

	year, month, day := time.Now().UTC().Date()
	nextMidnightTimePoint := time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC)
	dbGarbageCollector := database.NewGarbageCollector(db, logger, nextMidnightTimePoint, time.Hour*24)
	dbGarbageCollector.Run()

	repo := repository.NewRepository(db)
	serv := service.NewService(repo, logger)
	httpServer := v1.NewHttpServer(serv)

	// Starts server
	if err := httpServer.RunServer(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
		return
	}
}

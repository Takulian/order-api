package main

import (
	"context"
	"log"
	"net/http"
	"order-api/internal/cache"
	"order-api/internal/config"
	"order-api/internal/database"
	"order-api/internal/handler"
	"order-api/internal/observability"
	"order-api/internal/repository"
	"order-api/internal/router"
	"order-api/internal/service"

	_ "order-api/docs"
)

// @title Order API
// @version 1.0
// @description REST API belajar Go menggunakan net/http ServeMux
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Println("gagal load config")
		panic(err)
	}

	logger, shutdown, err := observability.InitLogging(ctx, "order-api", cfg.Telemetry)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			logger.Error("gagal shutdown logging", "error", err)
		}
	}()

	logger.Info("Logger berhasil di inisiasi")
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Error("gagal konek database", "error", err)
		panic(err)
	}
	defer db.Close()
	logger.Info("berhasil konek ke database")

	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		logger.Error("gagal konek redis", "error", err)
		panic(err)
	}
	defer rdb.Close()
	logger.Info("berhasil konek ke redis")

	repo := repository.NewPostgresRepository(db)
	cache := cache.NewRedisCache(rdb)
	service := service.NewOrderService(repo, cache, logger)
	orderHandler := handler.NewOrderHandler(service)
	router := router.NewRouter(orderHandler)
	logger.Info("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("server berhenti", "error", err)
	}
}

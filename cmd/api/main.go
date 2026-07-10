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
		log.Fatalf("gagal load config: %v", err)
	}

	logger, shutdown, err := observability.InitLogging(ctx, cfg.Telemetry)
	if err != nil {
		log.Fatalf("gagal setup logging: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Printf("gagal shutdown logging: %v", err)
		}
	}()

	logger.Info("order-api starting up")

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Error("koneksi ke database gagal", "error", err)
		panic(err)
	}
	defer db.Close()
	logger.Info("koneksi ke database berhasil")

	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		logger.Error("koneksi ke redis gagal", "error", err)
		panic(err)
	}
	defer rdb.Close()
	logger.Info("koneksi ke redis berhasil")

	repo := repository.NewPostgresRepository(db)
	cache := cache.NewRedisCache(rdb)
	service := service.NewOrderService(repo, cache)
	orderHandler := handler.NewOrderHandler(service)
	router := router.NewRouter(orderHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

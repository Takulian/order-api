package main

import (
	"log"
	"net/http"
	"order-api/internal/cache"
	"order-api/internal/config"
	"order-api/internal/database"
	"order-api/internal/handler"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
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

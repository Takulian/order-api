package main

import (
	"log"
	"net/http"
	"order-api/internal/cache"
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
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	redisClient, err := cache.NewRedis()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewPostgresRepository(db)
	service := service.NewOrderService(repo, redisClient)
	orderHandler := handler.NewOrderHandler(service)
	router := router.NewRouter(orderHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

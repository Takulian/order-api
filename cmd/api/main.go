package main

import(
	"log"
	"net/http"
	"order-api/internal/router"
	"order-api/internal/service"
	"order-api/internal/repository"
	"order-api/internal/handler"

	_ "order-api/docs"
)

// @title Order API
// @version 1.0
// @description REST API belajar Go menggunakan net/http ServeMux
// @host localhost:8080
// @BasePath /
func main() {
	repo := repository.NewOrderRepository()
	service := service.NewOrderService(repo)
	orderHandler := handler.NewOrderHandler(service)
	router := router.NewRouter(orderHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
package main

import(
	"log"
	"net/http"
	"order-api/internal/router"
	"order-api/internal/service"
	"order-api/internal/repository"
	"order-api/internal/handler"
)

func main() {
	repo := repository.NewOrderRepository()
	service := service.NewOrderService(repo)
	orderHandler := handler.NewOrderHandler(service)
	mux := router.NewRouter(orderHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
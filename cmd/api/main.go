package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"order-api/internal/cache"
	"order-api/internal/config"
	"order-api/internal/database"
	"order-api/internal/dto"
	"order-api/internal/event"
	"order-api/internal/grpcserver"
	"order-api/internal/handler"
	"order-api/internal/observability"
	"order-api/internal/repository"
	"order-api/internal/router"
	"order-api/internal/service"
	orderv1 "order-api/proto/order/v1"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	_ "order-api/docs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	rabbitConn, err := database.NewRabbitMQ(cfg.RabbitMQ.URL())
	if err != nil {
		logger.Error("gagal konek ke rabbitmq", "error", err)
	}
	defer rabbitConn.Close()

	publisher, err := event.NewRabbitMQPublisher(rabbitConn)
	if err != nil {
		logger.Error("publisher error", "error", err)
		panic(err)
	}
	defer publisher.Close()

	repo := repository.NewPostgresRepository(db)
	cache := cache.NewRedisCache(rdb)
	service := service.NewOrderService(repo, cache, publisher, logger)
	orderHandler := handler.NewOrderHandler(service, logger)
	router := router.NewRouter(orderHandler)
	srv := &http.Server{
		Addr:         cfg.App.PSN(),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	consumer, err := event.NewRabbitMQConsumer(rabbitConn)
	if err != nil {
		logger.Error("consumer error", "error", err)
		panic(err)
	}
	defer consumer.Close()

	go func() {
		err := consumer.ConsumeCheckout(ctx, func(ctx context.Context, evt event.CheckoutEvent) error {
			_, err := service.Create(ctx, dto.CreateOrderRequest{
				Customer: evt.Customer,
				Product:  evt.Product,
				Quantity: evt.Quantity,
			})
			return err
		})
		if err != nil {
			logger.Error("consumer order.checkout berhenti karena error", "error", err)
		}
	}()

	go func() {
		err := consumer.ConsumeOrderCreated(ctx, func(ctx context.Context, evt event.OrderCreatedEvent) error {
			logger.InfoContext(ctx, "menerima order.created",
				"order_id", evt.OrderID,
				"consumer", evt.Customer,
				"product", evt.Product,
			)
			return nil
		})
		if err != nil {
			logger.Error("consumer order.created berhenti karena error", "error", err)
		}
	}()

	go func() {
		logger.Info("starting server", "port", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("server berhenti", "error", err)
		}
	}()

	grpcOrderServer := grpcserver.NewOrderGRPCServer(service)
	grpcServer := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(grpcServer, grpcOrderServer)

	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			logger.Error("gRPC listener error", "error", err)
		}
		logger.Info("gRPC server listening on :9090")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC server error", "error", err)
		}
	}()

	gwMux := runtime.NewServeMux()
	if err := orderv1.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		gwMux,
		"localhost:9090",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		logger.Error("gateway registration error", "error", err)
	}

	logger.Info("gRPC-gateway listening on :8081")
	if err := http.ListenAndServe(":8081", gwMux); err != nil {
		logger.Error("gateway server error", "error", err)
	}
}

package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"loyalty_accounting/api"
	"loyalty_accounting/internal/config"
	"loyalty_accounting/internal/model"
	"loyalty_accounting/internal/service"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger *zap.Logger
	cfg    *config.Config
)

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
}

func main() {
	cfg = config.LoadConfig()
	err := runServer(cfg)
	if err != nil {
		logger.Fatal("Server initialization failed", zap.Error(err))
	}
}

func runServer(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.Username, cfg.Password, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	db.AutoMigrate(model.Account{}, model.Transaction{})

	service, err := service.NewService(logger, db)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.Host)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", cfg.Host, err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterLoyaltyAccountingServiceServer(grpcServer, service)

	go func() {
		logger.Info("gRPC server is running on port " + cfg.Host)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC server", zap.Error(err))
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	gracefulShutdown(grpcServer)
	return nil
}

func gracefulShutdown(grpcServer *grpc.Server) {
	logger.Info("Shutting down gRPC server gracefully...")
	grpcServer.GracefulStop()
	logger.Info("Server stopped gracefully")
}

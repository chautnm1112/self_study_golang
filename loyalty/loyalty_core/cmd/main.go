package main

import (
	_ "context"
	"flag"
	"fmt"
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
	apiLoyalty "github.com/chautnm1112/loyalty/loyalty_core/api"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/config"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/model"
	service "github.com/chautnm1112/loyalty/loyalty_core/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"syscall"
	_ "time"
)

var logger *zap.Logger
var cfg *config.Config

func run() error {
	cfg = config.LoadConfig()
	logger, _ = zap.NewProduction()
	err := serverAction(cfg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		logger.Error("Program start failed")
	}
}

func serverAction(cfg *config.Config) error {
	serverAddr := flag.String("addr", "loyalty-accounting:8112", "The server address in the format of host:port")

	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Fatal("Did not connect")
	}
	defer conn.Close()

	client := api.NewLoyaltyAccountingServiceClient(conn)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.Username, cfg.Password, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database")
	}

	db.AutoMigrate(model.Network{}, model.Merchant{}, model.Member{})

	logger, _ = zap.NewProduction()

	_, _ = service.NewService(logger, db, client)

	service, err := service.NewService(logger, db, client)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.Host)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", cfg.Host, err)
	}

	grpcServer := grpc.NewServer()
	apiLoyalty.RegisterLoyaltyCoreServiceServer(grpcServer, service)

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

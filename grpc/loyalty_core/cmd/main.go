package main

import (
	"context"
	_ "context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"loyalty_core/api"
	"loyalty_core/internal/config"
	"loyalty_core/internal/model"
	service "loyalty_core/internal/service"
	"time"
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

	logger.Info("cfg: " + cfg.Host)
	logger.Info("cfg: " + cfg.DBName)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	request := &api.CreateTransactionRequest{
		FromAccountId: 1,
		ToAccountId:   2,
		Point:         2000,
		Type:          api.TransactionType_REFUND_REDEEMED_POINTS,
	}
	//request := &api.EarnPointsRequest{
	//	Points:            1200,
	//	MemberAccountId:   1,
	//	MerchantAccountId: 2,
	//}
	//_, err = client.EarnPoints(ctx, request)
	_, err = client.CreateTransaction(ctx, request)
	if err != nil {
		logger.Fatal("Failed to add points", zap.Error(err))
		return nil
	}
	logger.Info("Create Transaction success")
	return nil
}

package repository

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"loyalty_accounting/internal/model"
	"os"
	"testing"
	"time"
)

var (
	gormDb *gorm.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "123321",
			"POSTGRES_DB":       "loyalty_accounting",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	defer postgresC.Terminate(ctx)

	_, err = postgresC.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	host, err := postgresC.Host(ctx)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, "postgres", "123321", "loyalty_accounting")
	gormDb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = gormDb.AutoMigrate(&model.Transaction{})
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

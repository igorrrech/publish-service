package main

import (
	"context"
	"log"
	"log/slog"
	"net"

	"github.com/igorrrech/publish-service/authorization/config"
	l "github.com/igorrrech/publish-service/authorization/pkg/logger"
	"github.com/igorrrech/publish-service/authorization/repo"
	"github.com/igorrrech/publish-service/authorization/service"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func main() {
	cfg := config.MustLoadConfig("./config.json")

	logger := slog.Default()

	userRepo := repo.NewUserRepository(cfg.UserDB.DsnString, logger)
	tokenRepo := repo.NewTokenRepository(cfg.Secret, cfg.AccessTTL)

	switch cfg.AppEnviroment {
	case "prod":
		logger = l.SetupLogger(l.PROD, cfg.Filepath)
	case "dev":
		logger = l.SetupLogger(l.DEV, cfg.Filepath)
	}

	srv := service.NewAuthService(
		net.JoinHostPort(cfg.Host, cfg.Port),
		*userRepo,
		*tokenRepo,
		cfg.AccessTTL,
		cfg.RefreshTTL,
		logger,
	)

	ctx := context.Background()
	srv.Run(ctx, cfg.ShutdownTimeout)
}

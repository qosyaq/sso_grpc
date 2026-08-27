// ! Приложение сервиса (для запуска различных компонентов)

package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"sso/internal/validator"
	"time"
)

type App struct {
	GRPCServ *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	databaseURL string,
	tokenTTL time.Duration,
) *App {
	storage, err := postgres.New(databaseURL)
	if err != nil {
		panic(err)
	}

	// Сервисный слой
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	validator := validator.New(log)

	grpcApp := grpcapp.New(log, authService, grpcPort, validator)

	return &App{
		GRPCServ: grpcApp,
	}
}

package app

import (
	"fmt"
	"log/slog"

	grpcapp "github.com/drobyshevv/classifier-expert-search/internal/app/grpc"
	"github.com/drobyshevv/classifier-expert-search/internal/services/expert"
	"github.com/drobyshevv/classifier-expert-search/internal/services/ml"
	"github.com/drobyshevv/classifier-expert-search/internal/storage/sqlite"
)

type App struct {
	GRPCSrv  *grpcapp.App
	MLClient *ml.MLServiceClient
}

func New(log *slog.Logger, grpcPort int, storagePath string) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// Подключаемся к AI Agent
	mlClient, err := ml.NewMLServiceClient("localhost:50052") // Порт твоего AI Agent
	if err != nil {
		panic(fmt.Sprintf("failed to connect to AI agent: %v", err))
	}

	// Сервисный слой
	expertService := expert.New(log, storage, storage, storage, storage, storage, storage, mlClient)

	grpcApp := grpcapp.New(log, expertService, grpcPort)

	return &App{
		GRPCSrv:  grpcApp,
		MLClient: mlClient,
	}
}

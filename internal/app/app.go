package app

import (
	"log/slog"

	grpcapp "github.com/drobyshevv/classifier-expert-search/internal/app/grpc"
	"github.com/drobyshevv/classifier-expert-search/internal/services/expert"
	"github.com/drobyshevv/classifier-expert-search/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string) *App {
	storage, err := sqlite.New(storagePath) // repository
	if err != nil {
		panic(err)
	}
	//Сервисный слой
	expertService := expert.New(log, storage, storage, storage, storage, storage, storage)

	grpcApp := grpcapp.New(log, expertService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}

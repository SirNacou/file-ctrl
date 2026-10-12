package main

import (
	"log"
	"log/slog"

	"github.com/sirnacou/file-ctrl/backend/internal/api"
	"github.com/sirnacou/file-ctrl/backend/internal/common"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage"
)

func main() {
	commonModule, err := common.NewCommonModule()
	if err != nil {
		log.Fatalln(err)
	}

	slog.SetDefault(commonModule.Logger)

	storageModule, err := storage.NewStorageModule(commonModule.Env.StorageDir, commonModule.DB)
	if err != nil {
		log.Fatalf("Failed to create Storage Module. Error: %e", err)
	}

	httpRouter := api.NewHttpRouter()
	httpRouter.RegisterRoutes(storageModule.RegisterAPI)

	if err := httpRouter.Run(commonModule.Env.Port()); err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v5"
	"github.com/sirnacou/file-ctrl/backend/internal/common/config"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load config. Error: %e", err)
	}

	router := echo.New()

	humaConfig := huma.DefaultConfig("File Ctrl API", "v1")
	humaConfig.DocsRenderer = huma.DocsRendererScalar
	api := humaecho.New(router, humaConfig)

	apiGrp := huma.NewGroup(api, "/api")

	storageModule, err := storage.NewStorageModule(cfg)
	if err != nil {
		log.Fatalf("Failed to create Storage Module. Error: %e", err)
	}

	storageModule.RegisterAPI(apiGrp)

	slog.Info("Server run on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", router))
}

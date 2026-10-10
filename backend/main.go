package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v5"
	"github.com/sirnacou/file-ctrl/backend/internal/common"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	commonModule, err := common.NewCommonModule(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	slog.SetDefault(commonModule.Logger)

	storageModule, err := storage.NewStorageModule(commonModule.Env)
	if err != nil {
		log.Fatalf("Failed to create Storage Module. Error: %e", err)
	}

	router := echo.New()

	humaConfig := huma.DefaultConfig("File Ctrl API", "v1")
	humaConfig.DocsRenderer = huma.DocsRendererScalar
	api := humaecho.New(router, humaConfig)

	apiGrp := huma.NewGroup(api, "/api")

	storageModule.RegisterAPI(apiGrp)

	slog.Info("Server run on port", slog.Int("PORT", int(commonModule.Env.Port())))
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", commonModule.Env.Port()), router))
}

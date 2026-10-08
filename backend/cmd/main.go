package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v5"
	"github.com/sirnacou/file-ctrl/backend/internal/features/storage"
)

func main() {
	router := echo.New()

	humaConfig := huma.DefaultConfig("File Ctrl API", "v1")
	humaConfig.DocsRenderer = huma.DocsRendererScalar
	api := humaecho.New(router, humaConfig)

	apiGrp := huma.NewGroup(api, "/api")

	storage.RegisterModule(apiGrp)

	slog.Info("Server run on port 8080")

	log.Fatalln(http.ListenAndServe("0.0.0.0:8080", router))
}

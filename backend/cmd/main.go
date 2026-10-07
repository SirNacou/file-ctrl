package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v5"
)

func main() {
	router := echo.New()

	humaConfig := huma.DefaultConfig("File Ctrl API", "v1")
	humaConfig.DocsRenderer = huma.DocsRendererScalar
	api := humaecho.New(router, humaConfig)

	apiGrp := huma.NewGroup(api, "/api")

	huma.Get(apiGrp, "", func(ctx context.Context, i *struct{}) (*struct {
		Status int
	}, error) {
		return &struct {
			Status int
		}{
			Status: http.StatusOK,
		}, nil
	})

	slog.Info("Server run on port 8080")

	log.Fatalln(http.ListenAndServe("0.0.0.0:8080", router))
}

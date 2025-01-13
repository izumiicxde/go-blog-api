package main

import (
	"log"
	"log/slog"

	"github.com/izumii.cxde/blog-api/cmd/api"
	"github.com/izumii.cxde/blog-api/config"
	"github.com/izumii.cxde/blog-api/storage"
)

func main() {
	db, err := storage.NewPostgresStorage(config.Envs)
	if err != nil {
		log.Fatal(err)
		return
	}
	server := api.NewAPIServer(config.Envs.Port, db)
	if err := server.Run(); err != nil {
		slog.Error("failed to run the server: ", slog.String("error", err.Error()))
	}
}

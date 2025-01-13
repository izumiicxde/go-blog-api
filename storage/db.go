package storage

import (
	"log/slog"

	"github.com/izumii.cxde/blog-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresStorage(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBAddress), &gorm.Config{})
	if err != nil {
		slog.Error("failed to open database: ", slog.String("error", err.Error()))
		return nil, err
	}
	slog.Info("database opened successfully")
	return db, nil
}

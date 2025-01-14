package storage

import (
	"log/slog"

	"github.com/izumii.cxde/blog-api/config"
	"github.com/izumii.cxde/blog-api/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresStorage(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBAddress), &gorm.Config{})
	if err != nil {
		slog.Error("failed to open database: ", slog.String("error", err.Error()))
		return nil, err
	}

	if err = db.AutoMigrate(
		&types.User{},
		&types.Blog{},
		&types.Tag{},
		&types.BlogTag{}); err != nil {
		slog.Error("failed to auto migrate: ", slog.String("error", err.Error()))
		return db, err
	} else {
		slog.Info("database auto migrated successfully")
	}

	slog.Info("database opened successfully")
	return db, nil
}

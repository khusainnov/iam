package infra

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/khusainnov/iam/app/config"
	"gitlab.com/khusainnov/driver/postgres"
	"go.uber.org/zap"
)

type ClientImpl struct {
	*sqlx.DB
}

func New(log *zap.Logger, cfg *config.DB) (*ClientImpl, error) {
	db, err := postgres.NewPostgresDB(
		postgres.ConfigPG{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			Password: cfg.Password,
			DBName:   cfg.Name,
			SSLMode:  cfg.SSLMode,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database, %w", err)
	}

	go func() {
		t := time.NewTicker(cfg.PingInterval)

		for range t.C {
			if err := db.Ping(); err != nil {
				log.Warn("failed to ping db", zap.Error(err))
			}
		}
	}()

	return &ClientImpl{db}, nil
}

func (db *ClientImpl) GetDB() *sqlx.DB {
	return db.DB
}


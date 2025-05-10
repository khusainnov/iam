package app

import (
	"github.com/khusainnov/iam/app/config"
	"github.com/khusainnov/iam/app/infra/log"
	"github.com/khusainnov/iam/app/service"
	"go.uber.org/zap"
)

type Container struct {
	Config   *config.Config
	Logger   *log.LoggerCtx
	Services *service.Services
	//Prometheus *prometheus.Registry
}

func NewContainer(logger *zap.Logger, cfg *config.Config) *Container {
	return &Container{
		Config:   cfg,
		Logger:   log.New(logger),
		Services: service.NewServices(),
	}
}
package app

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/khusainnov/iam/app/config"
	"github.com/khusainnov/iam/app/handler"
	"github.com/khusainnov/iam/app/infra/server"
	//"github.com/khusainnov/iam/app/infra/storage" 
	"github.com/khusainnov/iam/app/processor"
	//"github.com/khusainnov/iam/app/repository"
)

// App struct example
type App struct {
	Cfg    *config.Config
	Log    *zap.Logger
	Server *server.Server
	done   chan struct{}
}

func New() *App {
	cfg := config.NewFromEnv()

	app := &App{
		Cfg: cfg,
	}

	log, err := app.CreateLogger()
	if err != nil {
		panic(err)
	}

	app.Log = log

	c := NewContainer(log, cfg)

	app.Log.Info("config init", zap.Any("config", c.Config))

	app.Log.Info("connecting to db")
	/*conn, err := storage.New(app.Log, c.Config.DB)
	if err != nil {
		app.Log.Error("failed to connect to db")
		panic(err)
	}*/

	app.Log.Info("creating new handler")
	handlers := handler.New(
		c.Logger,
		processor.New(c.Logger),
	)

	app.Log.Info("creating new server")
	srv := server.New(c.Config.Server)
	app.Server = srv
	if err := app.Server.Init(handlers); err != nil {
		panic(err)
	}

	return app
}

func (a *App) CreateLogger() (*zap.Logger, error) {
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logCfg.EncoderConfig.TimeKey = "timestamp"

	return logCfg.Build(zap.AddCaller())
}

func (a *App) Run() {
	if a.Server != nil {
		go func() {
			defer func() { a.done <- struct{}{} }()
			a.Log.Info("starting server")
			err := a.Server.Run()
			if err != nil {
				a.Log.Error(err.Error(), zap.Error(err))
			}
		}()
	}

	a.Log.Info("iam running!")
	<-a.done
	a.stop()
}

func (a *App) stop() {
	a.Log.Info("iam stopping...")

	if err := a.Server.Stop(context.Background()); err != nil {
		a.Log.Error(err.Error(), zap.Error(err))
	}
}

package handler

import (
	"context"

	"github.com/khusainnov/iam/app/infra/log"
	"github.com/khusainnov/iam/app/processor"
	reqParams "github.com/khusainnov/iam/specs/server/params/model"
	respParams "github.com/khusainnov/iam/specs/server/response/model"
)

type API interface {
	EchoAPI
}

type EchoAPI interface {
	Echo(ctx context.Context, params *reqParams.EchoParams, resp *respParams.EchoResponse) error
}

// processor interface example
type Processor interface {
	Echo(ctx context.Context) error
}

// Handler example
type Handler struct {
	proc Processor
	logger *log.LoggerCtx
}

func New(logger *log.LoggerCtx, proc *processor.Processor) *Handler {
	return &Handler{
		logger: logger,
		proc: proc,
	}
}

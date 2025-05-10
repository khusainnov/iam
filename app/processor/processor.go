package processor

import (
	"context"

	"github.com/khusainnov/iam/app/infra/log"
)

// Processor example
type Processor struct {
	logCtx *log.LoggerCtx
}

func New(logCtx *log.LoggerCtx) *Processor {
	return &Processor{
		logCtx: logCtx,
	}
}

func (p *Processor) Echo(ctx context.Context) error {
	p.logCtx.Warn(ctx, "implement  me")

	return nil
}


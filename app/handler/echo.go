package handler

import (
	"net/http"

	"go.uber.org/zap"

	reqParams "github.com/khusainnov/iam/specs/server/params/model"
	respParams "github.com/khusainnov/iam/specs/server/response/model"
)

func (h *Handler) Echo(r *http.Request, params *reqParams.EchoParams, resp *respParams.EchoResponse) error {
	ctx := r.Context()
	
	h.logger.Info(ctx, "new call", zap.String("method", "echo"))

	resp.Message = "echo: " + params.Message

	return nil
}

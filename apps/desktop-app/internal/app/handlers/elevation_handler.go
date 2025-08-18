package handlers

import (
	"context"

	"github.com/oLenador/mulltbost/internal/app/container"

)

type BoosterHandler struct {
	ctx       context.Context
	container *container.Container
}

func NewBoosterHandler(container *container.Container) *BoosterHandler {
	return &BoosterHandler{
		container: container,
	}
}

func (h *BoosterHandler) SetContext(ctx context.Context) {
	h.ctx = ctx
}

func (h *ElevationHandler) checkIsAdmin
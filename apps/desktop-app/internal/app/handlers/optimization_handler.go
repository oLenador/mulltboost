package handlers

import (
	"context"

	"github.com/oLenador/mulltbost/internal/app/container"
	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
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

// Métodos expostos para o frontend via Wails

func (h *BoosterHandler) GetAvailableBoosters() []dto.BoosterDto {
	return h.container.BoosterService.GetAvailableBoosters()
}

func (h *BoosterHandler) GetBoosterState(id string) (*entities.BoosterState, error) {
	return h.container.BoosterService.GetBoosterState(id)
}

func (h *BoosterHandler) ApplyBooster(id string) (*entities.BoosterResult, error) {
	return h.container.BoosterService.ApplyBooster(h.ctx, id)
}

func (h *BoosterHandler) RevertBooster(id string) (*entities.BoosterResult, error) {
	return h.container.BoosterService.RevertBooster(h.ctx, id)
}

func (h *BoosterHandler) ApplyBoosterBatch(ids []string) (*entities.BatchResult, error) {
	return h.container.BoosterService.ApplyBoosterBatch(h.ctx, ids)
}

func (h *BoosterHandler) GetBoostersByCategory(category string) []entities.Booster {
	cat := entities.BoosterCategory(category)
	return h.container.BoosterService.GetBoostersByCategory(cat)
}

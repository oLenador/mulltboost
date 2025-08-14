package handlers

import (
	"context"
	"fmt"

	"github.com/oLenador/mulltbost/internal/app/container"
	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
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

func (h *BoosterHandler) GetAvailableBoosters(lang i18n.Language) []dto.BoosterDto {
	return h.container.BoosterService.GetAvailableBoosters(lang)
}

func (h *BoosterHandler) GetBoosterState(id string) (*entities.BoosterRollbackState, error) {
	return h.container.BoosterService.GetBoosterRollbackState(id)
}

func (h *BoosterHandler) ApplyBooster(id string) (*entities.BoosterResult, error) {
	return h.container.BoosterService.ApplyBooster(h.ctx, id)
}

func (h *BoosterHandler) RevertBooster(id string) (*entities.BoosterResult, error) {
	return h.container.BoosterService.RevertBooster(h.ctx, id)
}

func (h *BoosterHandler) ApplyBoosterBatch(ids []string) (*entities.BatchResult, error) {
	fmt.Print(ids)
	return h.container.BoosterService.ApplyBoosterBatch(h.ctx, ids)
}

func (h *BoosterHandler) GetBoostersByCategory(category entities.BoosterCategory, lang i18n.Language) []dto.BoosterDto {
	return h.container.BoosterService.GetBoostersByCategory(category, lang)
}

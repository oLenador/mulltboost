package handlers

import (
	"context"

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


func (h *BoosterHandler) GetOperationsHistory(id string) (*[]entities.BoostOperation, error) {
	return h.container.BoosterService.GetOperationsHistory(h.ctx, id)
}


func (h *BoosterHandler) GetAvailableBoosters(lang i18n.Language) []dto.BoosterDto {
	return h.container.BoosterService.GetAvailableBoosters(h.ctx, lang)
}
func (h *BoosterHandler) GetBoostersByCategory(category entities.BoosterCategory, lang i18n.Language) []dto.BoosterDto {
	return h.container.BoosterService.GetBoostersByCategory(h.ctx, category, lang)
}
func (h *BoosterHandler) GetExecutionQueueState() *entities.QueueState {
	return h.container.BoosterService.GetExecutionQueueState(h.ctx)
}


func (h *BoosterHandler) InitBoosterApply(id string) (entities.InitResult, error) {
    return h.container.BoosterService.InitBoosterApply(h.ctx, id)
}
func (h *BoosterHandler) InitBoosterApplyBatch(ids []string) (entities.InitResult, error) {
    return h.container.BoosterService.InitBoosterApplyBatch(h.ctx, ids)
}
func (h *BoosterHandler) InitRevertBooster(id string) (entities.InitResult, error) {
	return h.container.BoosterService.InitRevertBooster(h.ctx, id)
}
func (h *BoosterHandler) InitRevertBoosterBatch(ids []string) (entities.InitResult, error) {
	return h.container.BoosterService.InitRevertBoosterBatch(h.ctx, ids)
}


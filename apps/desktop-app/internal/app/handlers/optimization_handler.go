package handlers

import (
    "context"

    "github.com/oLenador/mulltbost/internal/app/container"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type OptimizationHandler struct {
    ctx       context.Context
    container *container.Container
}

func NewOptimizationHandler(container *container.Container) *OptimizationHandler {
    return &OptimizationHandler{
        container: container,
    }
}

func (h *OptimizationHandler) SetContext(ctx context.Context) {
    h.ctx = ctx
}

// Métodos expostos para o frontend via Wails

func (h *OptimizationHandler) GetAvailableOptimizations() []entities.Optimization {
    return h.container.OptimizationService.GetAvailableOptimizations()
}

func (h *OptimizationHandler) GetOptimizationState(id string) (*entities.OptimizationState, error) {
    return h.container.OptimizationService.GetOptimizationState(id)
}

func (h *OptimizationHandler) ApplyOptimization(id string) (*entities.OptimizationResult, error) {
    return h.container.OptimizationService.ApplyOptimization(h.ctx, id)
}

func (h *OptimizationHandler) RevertOptimization(id string) (*entities.OptimizationResult, error) {
    return h.container.OptimizationService.RevertOptimization(h.ctx, id)
}

func (h *OptimizationHandler) ApplyOptimizationBatch(ids []string) (*entities.BatchResult, error) {
    return h.container.OptimizationService.ApplyOptimizationBatch(h.ctx, ids)
}

func (h *OptimizationHandler) GetOptimizationsByCategory(category string) []entities.Optimization {
    cat := entities.OptimizationCategory(category)
    return h.container.OptimizationService.GetOptimizationsByCategory(cat)
}
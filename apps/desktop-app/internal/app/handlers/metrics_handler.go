package handlers

import (
    "context"

    "github.com/oLenador/mulltbost/internal/app/container"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type MetricsHandler struct {
    ctx       context.Context
    container *container.Container
}

func NewMetricsHandler(container *container.Container) *MetricsHandler {
    return &MetricsHandler{
        container: container,
    }
}

func (h *MetricsHandler) SetContext(ctx context.Context) {
    h.ctx = ctx
}

func (h *MetricsHandler) GetSystemMetrics() (*entities.SystemMetrics, error) {
    return h.container.MetricsService.GetSystemMetrics(h.ctx)
}

func (h *MetricsHandler) StartRealTimeMetrics(intervalSeconds int) error {
    return h.container.MetricsService.StartRealTimeMonitoring(h.ctx, intervalSeconds)
}

func (h *MetricsHandler) StopRealTimeMetrics() error {
    return h.container.MetricsService.StopRealTimeMonitoring(h.ctx)
}

func (h *MetricsHandler) IsMetrics() bool {
    return h.container.MetricsService.IsMonitoring()
}
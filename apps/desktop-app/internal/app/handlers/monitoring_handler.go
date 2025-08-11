package handlers

import (
    "context"

    "github.com/oLenador/mulltbost/internal/app/container"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type MonitoringHandler struct {
    ctx       context.Context
    container *container.Container
}

func NewMonitoringHandler(container *container.Container) *MonitoringHandler {
    return &MonitoringHandler{
        container: container,
    }
}

func (h *MonitoringHandler) SetContext(ctx context.Context) {
    h.ctx = ctx
}

func (h *MonitoringHandler) GetSystemMetrics() (*entities.SystemMetrics, error) {
    return h.container.MonitoringService.GetSystemMetrics(h.ctx)
}

func (h *MonitoringHandler) StartRealTimeMonitoring(intervalSeconds int) error {
    return h.container.MonitoringService.StartRealTimeMonitoring(h.ctx, intervalSeconds)
}

func (h *MonitoringHandler) StopRealTimeMonitoring() error {
    return h.container.MonitoringService.StopRealTimeMonitoring(h.ctx)
}

func (h *MonitoringHandler) IsMonitoring() bool {
    return h.container.MonitoringService.IsMonitoring()
}
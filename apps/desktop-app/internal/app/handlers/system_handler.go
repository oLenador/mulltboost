package handlers

import (
    "context"

    "github.com/oLenador/mulltbost/internal/app/container"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type SystemHandler struct {
    ctx       context.Context
    container *container.Container
}

func NewSystemHandler(container *container.Container) *SystemHandler {
    return &SystemHandler{
        container: container,
    }
}

func (h *SystemHandler) SetContext(ctx context.Context) {
    h.ctx = ctx
}

func (h *SystemHandler) GetSystemInfo() (*entities.SystemInfo, error) {
    return h.container.SystemInfoService.GetSystemInfo(h.ctx)
}

func (h *SystemHandler) GetHardwareInfo() (*entities.SystemInfo, error) {
    return h.container.SystemInfoService.GetHardwareInfo(h.ctx)
}

func (h *SystemHandler) RefreshSystemInfo() error {
    return h.container.SystemInfoService.RefreshSystemInfo(h.ctx)
}
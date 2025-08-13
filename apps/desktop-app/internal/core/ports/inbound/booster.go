package inbound

import (
	"context"

	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
)

type BoosterUseCase interface {
    Execute(ctx context.Context) (*entities.BoosterResult, error)
    Validate(ctx context.Context) error
    CanApply(ctx context.Context) bool
    CanRevert(ctx context.Context) bool
	GetEntity() entities.Booster 
	GetEntityDto(lang i18n.Language) dto.BoosterDto
    Revert(ctx context.Context) (*entities.BoosterResult, error)
}

type BoosterService interface {
    RegisterBooster(booster BoosterUseCase) error
    GetAvailableBoosters() []dto.BoosterDto
    GetBoosterState(id string) (*entities.BoosterState, error)
    ApplyBooster(ctx context.Context, id string) (*entities.BoosterResult, error)
    RevertBooster(ctx context.Context, id string) (*entities.BoosterResult, error)
    ApplyBoosterBatch(ctx context.Context, ids []string) (*entities.BatchResult, error)
    GetBoostersByCategory(category entities.BoosterCategory) []entities.Booster
}

type MonitoringService interface {
    GetSystemMetrics(ctx context.Context) (*entities.SystemMetrics, error)
    StartRealTimeMonitoring(ctx context.Context, interval int) error
    StopRealTimeMonitoring(ctx context.Context) error
    IsMonitoring() bool
}

type SystemInfoService interface {
    GetSystemInfo(ctx context.Context) (*entities.SystemInfo, error)
    GetHardwareInfo(ctx context.Context) (*entities.SystemInfo, error)
    RefreshSystemInfo(ctx context.Context) error
}

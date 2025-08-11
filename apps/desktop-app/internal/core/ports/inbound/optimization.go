package inbound

import (
    "context"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type OptimizationUseCase interface {
    Execute(ctx context.Context) (*entities.OptimizationResult, error)
    Validate(ctx context.Context) error
    CanApply(ctx context.Context) bool
    CanRevert(ctx context.Context) bool
    GetInfo() entities.Optimization
    Revert(ctx context.Context) (*entities.OptimizationResult, error)
}

type OptimizationService interface {
    GetAvailableOptimizations() []entities.Optimization
    GetOptimizationState(id string) (*entities.OptimizationState, error)
    ApplyOptimization(ctx context.Context, id string) (*entities.OptimizationResult, error)
    RevertOptimization(ctx context.Context, id string) (*entities.OptimizationResult, error)
    ApplyOptimizationBatch(ctx context.Context, ids []string) (*entities.BatchResult, error)
    GetOptimizationsByCategory(category entities.OptimizationCategory) []entities.Optimization
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
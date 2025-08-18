package inbound

import (
	"context"

	windows "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
)

type BoosterUseCase interface {
    Execute(ctx context.Context) (*entities.BoostApplyResult, error)
    Validate(ctx context.Context) error
    CanApply(ctx context.Context) bool
    CanRevert(ctx context.Context) bool
	GetEntity() entities.Booster 
	GetEntityDto(lang i18n.Language) dto.BoosterDto
    Revert(ctx context.Context) (*entities.BoostRevertResult, error)
}

type BoosterService interface {
	RegisterBooster(booster BoosterUseCase) error

	GetOperationsHistory(ctx context.Context, id string) (*[]entities.BoostOperation, error)
	GetAvailableBoosters(ctx context.Context, lang i18n.Language) []dto.BoosterDto
	GetBoosterStatus(ctx context.Context, id string) (*dto.BoosterDto, error)
	GetBoostersByCategory(ctx context.Context, category entities.BoosterCategory, lang i18n.Language) []dto.BoosterDto
	GetExecutionQueueState(ctx context.Context) (*[]dto.BoosterDto, error)
	InitBoosterApply(ctx context.Context, id string) (entities.AsyncOperationResult, error)
	InitBoosterApplyBatch(ctx context.Context, ids []string) (entities.AsyncOperationResult, error)
	InitRevertBooster(ctx context.Context, id string) (*entities.AsyncOperationResult, error)
	InitRevertBoosterBatch(ctx context.Context, ids []string) (*entities.AsyncOperationResult, error)
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


type PlatformExecutor interface {
	Execute(ctx context.Context, boosterID string) (*entities.BoostApplyResult, error)
	Validate(ctx context.Context) error
	CanExecute(ctx context.Context) bool
	Revert(ctx context.Context, backupData entities.BackupData) (*entities.BoostRevertResult, error)
}

type ExecutorFactory interface {
	CreateExecutor(boosterID string) (PlatformExecutor, error)
}

type ExecutorDepServices struct {
	TcpService       windows.TCPOptimizationService
	RegistryService  windows.RegistryService
	SystemService    windows.SystemAPIService
	ElevationService windows.ElevationService
	NetworkService   windows.NetworkAdapterService
	MemoryService    windows.MemoryManagementService
	GpuService       windows.GPUOptimizationService
	PowerService     windows.PowerManagementService
	ProcessService   windows.ProcessPriorityService
	CacheService     windows.CacheManagementService
	FileService      windows.FileSystemService
}
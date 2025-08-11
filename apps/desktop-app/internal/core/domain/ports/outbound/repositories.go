package outbound

import (
    "context"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type OptimizationStateRepository interface {
    Save(ctx context.Context, state *entities.OptimizationState) error
    GetByID(ctx context.Context, id string) (*entities.OptimizationState, error)
    GetAll(ctx context.Context) ([]*entities.OptimizationState, error)
    Delete(ctx context.Context, id string) error
}

type SystemMetricsRepository interface {
    GetCPUMetrics(ctx context.Context) (*entities.CPUMetrics, error)
    GetMemoryMetrics(ctx context.Context) (*entities.MemoryMetrics, error)
    GetGPUMetrics(ctx context.Context) (*entities.GPUMetrics, error)
    GetNetworkMetrics(ctx context.Context) (*entities.NetworkMetrics, error)
    GetTemperatureMetrics(ctx context.Context) (*entities.TemperatureMetrics, error)
    GetDiskMetrics(ctx context.Context) (*entities.DiskMetrics, error)
}

type SystemInfoRepository interface {
    GetOSInfo(ctx context.Context) (*entities.OSInfo, error)
    GetCPUInfo(ctx context.Context) (*entities.CPUInfo, error)
    GetMemoryInfo(ctx context.Context) (*entities.MemoryInfo, error)
    GetGPUInfo(ctx context.Context) ([]entities.GPUInfo, error)
    GetStorageInfo(ctx context.Context) ([]entities.StorageInfo, error)
    GetNetworkInfo(ctx context.Context) ([]entities.NetworkInfo, error)
    GetMotherboardInfo(ctx context.Context) (*entities.MotherboardInfo, error)
}

type RegistryRepository interface {
    GetValue(ctx context.Context, key, name string) (interface{}, error)
    SetValue(ctx context.Context, key, name string, value interface{}) error
    BackupKey(ctx context.Context, key string) (map[string]interface{}, error)
    RestoreKey(ctx context.Context, key string, backup map[string]interface{}) error
}

type ServiceRepository interface {
    GetServiceState(ctx context.Context, name string) (string, error)
    SetServiceState(ctx context.Context, name string, state string) error
    ListServices(ctx context.Context) ([]string, error)
}
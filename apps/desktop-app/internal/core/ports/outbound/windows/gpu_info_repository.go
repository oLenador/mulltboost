package windows

import (
	"context"
)

type GPUInfoRepository interface {
	// Informações da GPU
	GetGPUInfo(ctx context.Context) ([]*entities.GPUInfo, error)
	GetPrimaryGPU(ctx context.Context) (*entities.GPUInfo, error)
	RefreshGPUInfo(ctx context.Context) error
	
	// Status
	GetGPUUsage(ctx context.Context, gpuID string) (float64, error)
	GetGPUTemperature(ctx context.Context, gpuID string) (float64, error)
	GetVRAMUsage(ctx context.Context, gpuID string) (float64, error)
}
package windows

import (
	"context"
)

type GPUOptimizationService interface {
	// Otimização
	OptimizeGPUForGaming(ctx context.Context) error
	OptimizeGPUForWorkstation(ctx context.Context) error
	RestoreGPUDefaults(ctx context.Context) error
	
	// Configurações
	SetGPUPowerMode(ctx context.Context, mode string) error
	SetGPUMemoryClock(ctx context.Context, clockSpeed int) error
	SetGPUCoreClock(ctx context.Context, clockSpeed int) error
	
	// Status
	GetGPUConfiguration(ctx context.Context) (*entities.GPUConfiguration, error)
	IsGPUOptimized(ctx context.Context) (bool, error)
}
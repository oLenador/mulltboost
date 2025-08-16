package windows

import (
	"context"
)

type PowerManagementService interface {
	// Perfis de energia
	SetPowerProfile(ctx context.Context, profileName string) error
	GetActivePowerProfile(ctx context.Context) (*entities.PowerProfile, error)
	GetAvailablePowerProfiles(ctx context.Context) ([]*entities.PowerProfile, error)
	
	// Configurações avançadas
	DisablePowerThrottling(ctx context.Context) error
	SetCPUPowerPolicy(ctx context.Context, policy string) error
	SetUSBPowerSettings(ctx context.Context, enabled bool) error
	
	// Criação de perfis
	CreateCustomPowerProfile(ctx context.Context, profile *entities.PowerProfile) error
	DeletePowerProfile(ctx context.Context, profileID string) error
}
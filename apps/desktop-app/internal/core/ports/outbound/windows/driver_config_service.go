package windows

import (
	"context"
)

type DriverConfigService interface {
	// Gerenciamento de drivers
	UpdateDrivers(ctx context.Context) error
	GetOutdatedDrivers(ctx context.Context) ([]*entities.DriverConfiguration, error)
	InstallDriver(ctx context.Context, driverPath string) error
	
	// Configuração
	OptimizeDriverSettings(ctx context.Context) error
	RestoreDriverDefaults(ctx context.Context) error
	
	// Status
	GetDriverStatus(ctx context.Context, deviceID string) (*entities.DriverConfiguration, error)
	ListInstalledDrivers(ctx context.Context) ([]*entities.DriverConfiguration, error)
}
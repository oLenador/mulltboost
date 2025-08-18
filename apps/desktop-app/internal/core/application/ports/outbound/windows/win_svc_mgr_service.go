package windows

import (
	"context"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type WinServiceaManagerService interface {
	// Gerenciamento de serviços
	StartService(ctx context.Context, serviceName string) error
	StopService(ctx context.Context, serviceName string) error
	RestartService(ctx context.Context, serviceName string) error
	
	// Configuração
	SetServiceStartupType(ctx context.Context, serviceName string, startupType string) error
	DisableUnnecessaryServices(ctx context.Context) error
	EnableEssentialServices(ctx context.Context) error
	
	// Status
	GetServiceStatus(ctx context.Context, serviceName string) (*entities.WindowsService, error)
	ListServices(ctx context.Context) ([]*entities.WindowsService, error)
	IsServiceRunning(ctx context.Context, serviceName string) (bool, error)
}
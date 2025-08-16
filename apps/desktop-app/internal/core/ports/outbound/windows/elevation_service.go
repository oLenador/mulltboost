package windows

import "context"

type ElevationService interface {
	// Verificação de privilégios
	IsElevated(ctx context.Context) (bool, error)
	RequiresElevation(ctx context.Context, operation string) (bool, error)
	
	// Elevação
	RequestElevation(ctx context.Context) error
	RunAsAdmin(ctx context.Context, command string, args []string) error
	
	// UAC
	IsUACEnabled(ctx context.Context) (bool, error)
	ConfigureUAC(ctx context.Context, enabled bool) error
}
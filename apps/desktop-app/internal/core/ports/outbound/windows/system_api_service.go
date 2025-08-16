package windows

import "context"

type SystemAPIService interface {
	// APIs do sistema
	CallWindowsAPI(ctx context.Context, apiName string, params []interface{}) (interface{}, error)
	LoadSystemDLL(ctx context.Context, dllName string) error
	UnloadSystemDLL(ctx context.Context, dllName string) error
	
	// Informações do sistema
	GetSystemVersion(ctx context.Context) (string, error)
	GetSystemArchitecture(ctx context.Context) (string, error)
	IsWindows11(ctx context.Context) (bool, error)
	
	// Recursos avançados
	ExecutePowerShellCommand(ctx context.Context, command string) (string, error)
	ExecuteCMDCommand(ctx context.Context, command string) (string, error)
	RunSystemCommand(ctx context.Context, command string, args []string) (string, error)
}

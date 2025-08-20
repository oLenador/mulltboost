//go:build windows
package connection

import (
	"context"
	"fmt"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	windows "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
)

type TCPFastOpenExecutor struct {
	registryService  windows.RegistryService
	systemService    windows.SystemAPIService
	elevationService windows.ElevationService
}

func NewTCPFastOpenExecutor(
	registryService windows.RegistryService,
	systemService windows.SystemAPIService,
	elevationService windows.ElevationService,
) *TCPFastOpenExecutor {
	return &TCPFastOpenExecutor{
		registryService:  registryService,
		systemService:    systemService,
		elevationService: elevationService,
	}
}

func (e *TCPFastOpenExecutor) Execute(ctx context.Context, boosterID string) (*entities.BoostApplyResult, error) {
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return &entities.BoostApplyResult{
			Success: false,
			Message: "Administrator privileges required for TCP Fast Open",
			Error:   fmt.Errorf("elevation required"),
		}, fmt.Errorf("elevation required")
	}

	backupData := make(map[string]interface{})

	// Verificar se é Windows 10 versão 1607 ou superior (suporte ao TCP Fast Open)
	isWin11, err := e.systemService.IsWindows11(ctx)
	if err != nil {
		return &entities.BoostApplyResult{
			Success: false,
			Message: "Failed to check Windows version",
			Error:   err,
		}, err
	}

	// TCP Fast Open é suportado no Windows 10 versão 1607+ e Windows 11
	registryPath := `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`
	
	// Backup configuração atual
	currentValue, err := e.registryService.ReadRegistryValue(ctx, registryPath, "TcpFastOpen")
	if err == nil && currentValue != nil {
		backupData["TcpFastOpen"] = currentValue
	}

	currentValueClient, err := e.registryService.ReadRegistryValue(ctx, registryPath, "TcpFastOpenClient")
	if err == nil && currentValueClient != nil {
		backupData["TcpFastOpenClient"] = currentValueClient
	}

	// Habilitar TCP Fast Open
	if err := e.registryService.WriteRegistryValue(ctx, registryPath, "TcpFastOpen", 1); err != nil {
		return &entities.BoostApplyResult{
			Success: false,
			Message: "Failed to enable TCP Fast Open",
			Error:   err,
			BackupData: backupData,
		}, err
	}

	// Habilitar TCP Fast Open para cliente (Windows 10 1703+)
	if isWin11 {
		if err := e.registryService.WriteRegistryValue(ctx, registryPath, "TcpFastOpenClient", 1); err != nil {
			return &entities.BoostApplyResult{
				Success: false,
				Message: "Failed to enable TCP Fast Open client mode",
				Error:   err,
				BackupData: backupData,
			}, err
		}
	}

	// Configurações adicionais para otimização de latência inicial
	additionalTweaks := map[string]interface{}{
		"TcpInitialRtt":      1000, // Reduzir RTT inicial
		"TcpMaxConnectRetransmissions": 2, // Reduzir tentativas de conexão
	}

	for valueName, value := range additionalTweaks {
		currentVal, _ := e.registryService.ReadRegistryValue(ctx, registryPath, valueName)
		if currentVal != nil {
			backupData[valueName] = currentVal
		}

		if err := e.registryService.WriteRegistryValue(ctx, registryPath, valueName, value); err != nil {
			return &entities.BoostApplyResult{
				Success: false,
				Message: fmt.Sprintf("Failed to set %s", valueName),
				Error:   err,
				BackupData: backupData,
			}, err
		}
	}

	return &entities.BoostApplyResult{
		Success:    true,
		Message:    "TCP Fast Open enabled successfully - system restart recommended",
		BackupData: backupData,
		Error:      nil,
	}, nil
}

func (e *TCPFastOpenExecutor) Validate(ctx context.Context) error {
	if e.registryService == nil {
		return fmt.Errorf("registry service not available")
	}

	// Verificar versão do Windows
	version, err := e.systemService.GetSystemVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get system version: %w", err)
	}

	if version == "" {
		return fmt.Errorf("unsupported operating system")
	}

	return nil
}

func (e *TCPFastOpenExecutor) CanExecute(ctx context.Context) bool {
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return false
	}

	// Verificar se TCP Fast Open já está habilitado
	registryPath := `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`
	current, err := e.registryService.ReadRegistryValue(ctx, registryPath, "TcpFastOpen")
	
	if err != nil {
		return true // Pode ser configurado
	}

	// Se já está habilitado, não precisa executar
	if current == 1 {
		return false
	}

	return true
}

func (e *TCPFastOpenExecutor) Revert(ctx context.Context, backupData entities.BackupData) (*entities.BoostRevertResult, error) {
	if backupData == nil {
		return &entities.BoostRevertResult{
			Success: false,
			Message: "No backup data available for revert",
			Error:   fmt.Errorf("backup data is nil"),
		}, fmt.Errorf("backup data is nil")
	}

	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return &entities.BoostRevertResult{
			Success: false,
			Message: "Administrator privileges required for TCP Fast Open revert",
			Error:   fmt.Errorf("elevation required"),
		}, fmt.Errorf("elevation required")
	}

	registryPath := `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`

	// Restaurar valores do backup
	for valueName, value := range backupData {
		if err := e.registryService.WriteRegistryValue(ctx, registryPath, valueName, value); err != nil {
			return &entities.BoostRevertResult{
				Success: false,
				Message: fmt.Sprintf("Failed to restore %s", valueName),
				Error:   err,
			}, err
		}
	}

	return &entities.BoostRevertResult{
		Success: true,
		Message: "TCP Fast Open reverted successfully - system restart recommended",
		Error:   nil,
	}, nil
}
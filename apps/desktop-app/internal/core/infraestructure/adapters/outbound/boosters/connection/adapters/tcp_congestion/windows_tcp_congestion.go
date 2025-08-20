//go:build windows
package connection

import (
	"context"
	"fmt"
	"strings"

	windows "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type TCPCongestionExecutor struct {
	tcpService       windows.TCPOptimizationService
	registryService  windows.RegistryService
	systemService    windows.SystemAPIService
	elevationService windows.ElevationService
}

func NewTCPCongestionExecutor(
	tcpService windows.TCPOptimizationService,
	registryService windows.RegistryService,
	systemService windows.SystemAPIService,
	elevationService windows.ElevationService,
) *TCPCongestionExecutor {
	return &TCPCongestionExecutor{
		tcpService:       tcpService,
		registryService:  registryService,
		systemService:    systemService,
		elevationService: elevationService,
	}
}

func (e *TCPCongestionExecutor) Execute(ctx context.Context, boosterID string) (*entities.BoostApplyResult, error) {
	// Verificar elevação
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return &entities.BoostApplyResult{
			Success: false,
			Message: "Administrator privileges required for TCP congestion control",
			Error:   fmt.Errorf("elevation required"),
		}, fmt.Errorf("elevation required")
	}

	backupData := make(map[string]interface{})

	// Backup configuração atual de congestion control
	currentAlgorithm, err := e.registryService.ReadRegistryValue(ctx, 
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, 
		"TcpCongestionControl")
	if err == nil && currentAlgorithm != nil {
		backupData["TcpCongestionControl"] = currentAlgorithm
	}

	// Verificar se é Windows 10/11 para usar CUBIC, senão usar CTCP
	isWin11, err := e.systemService.IsWindows11(ctx)
	algorithm := "ctcp" // Default para versões mais antigas
	if err == nil && isWin11 {
		algorithm = "cubic"
	}

	// Configurar algoritmo de congestionamento TCP
	if err := e.tcpService.SetTCPCongestionControl(ctx, algorithm); err != nil {
		return &entities.BoostApplyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to set TCP congestion control to %s", algorithm),
			Error:   err,
			BackupData: backupData,
		}, err
	}

	registryTweaks := map[string]interface{}{
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\TcpInitialRtt`:     3000,
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\TcpMaxDupAcks`:     2,
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\EnableWsd`:        0,
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Tcp1323Opts`:      3,
	}

	for keyPath, value := range registryTweaks {
		key := keyPath[:strings.LastIndex(keyPath, `\`)]
		valueName := keyPath[strings.LastIndex(keyPath, `\`)+1:]
		
		// Backup valor atual
		currentValue, _ := e.registryService.ReadRegistryValue(ctx, key, valueName)
		if currentValue != nil {
			backupData[keyPath] = currentValue
		}

		// Aplicar nova configuração
		if err := e.registryService.WriteRegistryValue(ctx, key, valueName, value); err != nil {
			return &entities.BoostApplyResult{
				Success: false,
				Message: fmt.Sprintf("Failed to write registry value: %s", keyPath),
				Error:   err,
				BackupData: backupData,
			}, err
		}
	}

	return &entities.BoostApplyResult{
		Success:    true,
		Message:    fmt.Sprintf("TCP congestion control set to %s successfully", strings.ToUpper(algorithm)),
		BackupData: backupData,
		Error:      nil,
	}, nil
}

func (e *TCPCongestionExecutor) Validate(ctx context.Context) error {
	if e.tcpService == nil {
		return fmt.Errorf("TCP optimization service not available")
	}

	if e.registryService == nil {
		return fmt.Errorf("registry service not available")
	}

	// Verificar se o sistema suporta controle de congestionamento
	version, err := e.systemService.GetSystemVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get system version: %w", err)
	}

	if version == "" {
		return fmt.Errorf("unsupported operating system")
	}

	return nil
}

func (e *TCPCongestionExecutor) CanExecute(ctx context.Context) bool {
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return false
	}

	// Verificar se já está configurado
	current, err := e.registryService.ReadRegistryValue(ctx,
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
		"TcpCongestionControl")
	
	if err != nil {
		return true // Pode ser configurado
	}

	// Se já está configurado como CUBIC ou CTCP, não precisa executar
	if current == "cubic" || current == "ctcp" {
		return false
	}

	return true
}

func (e *TCPCongestionExecutor) Revert(ctx context.Context, backupData entities.BackupData) (*entities.BoostRevertResult, error) {
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
			Message: "Administrator privileges required for TCP revert",
			Error:   fmt.Errorf("elevation required"),
		}, fmt.Errorf("elevation required")
	}

	// Restaurar valores do registro
	for keyPath, value := range backupData {
		key := keyPath[:strings.LastIndex(keyPath, `\`)]
		valueName := keyPath[strings.LastIndex(keyPath, `\`)+1:]
		
		if err := e.registryService.WriteRegistryValue(ctx, key, valueName, value); err != nil {
			return &entities.BoostRevertResult{
				Success: false,
				Message: fmt.Sprintf("Failed to restore registry value: %s", keyPath),
				Error:   err,
			}, err
		}
	}

	return &entities.BoostRevertResult{
		Success: true,
		Message: "TCP congestion control reverted successfully",
		Error:   nil,
	}, nil
}

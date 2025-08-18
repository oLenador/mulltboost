//go:build windows
package connection

import (
	"context"
	"fmt"

	windows "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type TCPAdvancedExecutor struct {
	tcpService    windows.TCPOptimizationService
	registryService windows.RegistryService
	systemService windows.SystemAPIService
	elevationService windows.ElevationService
}

func NewTCPAdvancedExecutor(
	tcpService windows.TCPOptimizationService,
	registryService windows.RegistryService,
	systemService windows.SystemAPIService,
	elevationService windows.ElevationService,
) *TCPAdvancedExecutor {
	return &TCPAdvancedExecutor{
		tcpService:    tcpService,
		registryService: registryService,
		systemService: systemService,
		elevationService: elevationService,
	}
}

func (e *TCPAdvancedExecutor) Execute(ctx context.Context, boosterID string) (*entities.BoosterResult, error) {
	// Verificar se precisa de elevação
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: "Failed to check elevation status",
			Error:   err,
		}, err
	}

	if !elevated {
		return &entities.BoosterResult{
			Success: false,
			Message: "Administrator privileges required for TCP optimization",
			Error:   fmt.Errorf("elevation required"),
		}, fmt.Errorf("elevation required")
	}

	backupData := make(map[string]interface{})

	// Backup configurações atuais
	currentConfig, err := e.tcpService.GetTCPConfiguration(ctx)
	if err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: "Failed to get current TCP configuration",
			Error:   err,
		}, err
	}
	backupData["tcp_config"] = currentConfig

	// Aplicar otimizações TCP avançadas
	// 1. Otimizar TCP para gaming
	if err := e.tcpService.OptimizeTCPForGaming(ctx); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: "Failed to optimize TCP for gaming",
			Error:   err,
			BackupData: backupData,
		}, err
	}

	// 2. Configurar tamanho da janela TCP otimizado
	if err := e.tcpService.SetTCPWindowSize(ctx, 65536); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: "Failed to set TCP window size",
			Error:   err,
			BackupData: backupData,
		}, err
	}

	// 3. Desabilitar algoritmo de Nagle para reduzir latência
	if err := e.tcpService.DisableNagleAlgorithm(ctx); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: "Failed to disable Nagle algorithm",
			Error:   err,
			BackupData: backupData,
		}, err
	}

	// 4. Configurações avançadas no registro
	registryTweaks := map[string]interface{}{
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\TcpAckFrequency`: 1,
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\TCPNoDelay`:      1,
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\TcpDelAckTicks`:  0,
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\MaxConnectionsPerServer`: 16,
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\MaxConnectionsPer1_0Server`: 16,
	}

	for keyPath, value := range registryTweaks {
		// Fazer backup do valor atual
		currentValue, _ := e.registryService.ReadRegistryValue(ctx, keyPath, "")
		if currentValue != nil {
			backupData[keyPath] = currentValue
		}

		// Aplicar nova configuração
		if err := e.registryService.WriteRegistryValue(ctx, keyPath, "", value); err != nil {
			return &entities.BoosterResult{
				Success: false,
				Message: fmt.Sprintf("Failed to write registry value: %s", keyPath),
				Error:   err,
				BackupData: backupData,
			}, err
		}
	}

	return &entities.BoosterResult{
		Success:    true,
		Message:    "Advanced TCP/IP optimization applied successfully",
		BackupData: backupData,
		Error:      nil,
	}, nil
}

func (e *TCPAdvancedExecutor) Validate(ctx context.Context) error {
	// Verificar se os serviços estão funcionando
	if e.tcpService == nil {
		return fmt.Errorf("TCP optimization service not available")
	}

	if e.registryService == nil {
		return fmt.Errorf("registry service not available")
	}

	// Verificar se é Windows
	version, err := e.systemService.GetSystemVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get system version: %w", err)
	}

	if version == "" {
		return fmt.Errorf("unsupported operating system")
	}

	return nil
}

func (e *TCPAdvancedExecutor) CanExecute(ctx context.Context) bool {
	// Verificar privilégios de administrador
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return false
	}

	// Verificar se já está otimizado
	optimized, err := e.tcpService.IsTCPOptimized(ctx)
	if err != nil {
		return false
	}

	return !optimized
}

func (e *TCPAdvancedExecutor) Revert(ctx context.Context, backupData entities.BackupData) (*entities.BoosterResult, error) {
	if backupData == nil {
		return &entities.BoosterResult{
			Success: false,
			Message: "No backup data available for revert",
			Error:   fmt.Errorf("backup data is nil"),
		}, fmt.Errorf("backup data is nil")
	}

	// Verificar elevação
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return &entities.BoosterResult{
			Success: false,
			Message: "Administrator privileges required for TCP revert",
			Error:   fmt.Errorf("elevation required"),
		}, fmt.Errorf("elevation required")
	}

	// Restaurar configurações TCP padrão
	if err := e.tcpService.RestoreDefaultTCPSettings(ctx); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: "Failed to restore default TCP settings",
			Error:   err,
		}, err
	}

	// Restaurar configurações do registro
	for key, value := range backupData {
		if key == "tcp_config" {
			continue // Já restaurado pelo método acima
		}

		if err := e.registryService.WriteRegistryValue(ctx, key, "", value); err != nil {
			return &entities.BoosterResult{
				Success: false,
				Message: fmt.Sprintf("Failed to restore registry value: %s", key),
				Error:   err,
			}, err
		}
	}

	return &entities.BoosterResult{
		Success: true,
		Message: "Advanced TCP/IP optimization reverted successfully",
		Error:   nil,
	}, nil
}
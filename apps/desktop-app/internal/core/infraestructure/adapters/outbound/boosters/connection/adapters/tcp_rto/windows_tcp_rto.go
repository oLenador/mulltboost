//go:build windows
package connection

import (
	"context"
	"fmt"

	windows "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type TCPRTOExecutor struct {
	registryService  windows.RegistryService
	systemService    windows.SystemAPIService
	elevationService windows.ElevationService
}

func NewTCPRTOExecutor(
	registryService windows.RegistryService,
	systemService windows.SystemAPIService,
	elevationService windows.ElevationService,
) *TCPRTOExecutor {
	return &TCPRTOExecutor{
		registryService:  registryService,
		systemService:    systemService,
		elevationService: elevationService,
	}
}

func (e *TCPRTOExecutor) Execute(ctx context.Context, boosterID string) (*entities.BoostApplyResult, error) {
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return &entities.BoostApplyResult{
			Success: false,
			Message: "Administrator privileges required for TCP RTO adjustment",
			Error:   fmt.Errorf("elevation required"),
		}, fmt.Errorf("elevation required")
	}

	backupData := make(map[string]interface{})
	registryPath := `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`

	// Configurações de RTO (Retransmission Timeout) para otimização
	rtoTweaks := map[string]interface{}{
		"TcpInitialRtt":           1000, // RTT inicial em ms (padrão: 3000)
		"TcpMaxDataRetransmissions": 3,   // Máximo de retransmissões (padrão: 5)
		"TcpMaxConnectRetransmissions": 2, // Retransmissões de conexão (padrão: 3)
		"TcpTimedWaitDelay":       30,   // Delay em TIME_WAIT (padrão: 240)
		"TcpFinWait2Timeout":      40,   // Timeout FIN_WAIT_2 (padrão: 240)
		"KeepAliveTime":           300000, // Keep-alive em ms (padrão: 7200000)
		"KeepAliveInterval":       1000,   // Intervalo keep-alive (padrão: 1000)
	}

	// Fazer backup dos valores atuais
	for valueName := range rtoTweaks {
		currentValue, err := e.registryService.ReadRegistryValue(ctx, registryPath, valueName)
		if err == nil && currentValue != nil {
			backupData[valueName] = currentValue
		}
	}

	// Aplicar as novas configurações
	for valueName, value := range rtoTweaks {
		if err := e.registryService.WriteRegistryValue(ctx, registryPath, valueName, value); err != nil {
			return &entities.BoostApplyResult{
				Success: false,
				Message: fmt.Sprintf("Failed to set %s", valueName),
				Error:   err,
				BackupData: backupData,
			}, err
		}
	}

	// Aplicar configurações globais de TCP que afetam RTO
	globalTweaks := map[string]interface{}{
		"TcpAckFrequency": 1, // Enviar ACK para cada pacote
		"TcpDelAckTicks":  0, // Desabilitar delayed ACK
		"TCPNoDelay":      1, // Desabilitar algoritmo de Nagle
	}

	for valueName, value := range globalTweaks {
		currentValue, err := e.registryService.ReadRegistryValue(ctx, registryPath, valueName)
		if err == nil && currentValue != nil {
			backupData[valueName] = currentValue
		}

		if err := e.registryService.WriteRegistryValue(ctx, registryPath, valueName, value); err != nil {
			return &entities.BoostApplyResult{
				Success: false,
				Message: fmt.Sprintf("Failed to set global TCP parameter %s", valueName),
				Error:   err,
				BackupData: backupData,
			}, err
		}
	}

	return &entities.BoostApplyResult{
		Success:    true,
		Message:    "TCP retransmission interval optimized successfully",
		BackupData: backupData,
		Error:      nil,
	}, nil
}

func (e *TCPRTOExecutor) Validate(ctx context.Context) error {
	if e.registryService == nil {
		return fmt.Errorf("registry service not available")
	}

	if e.systemService == nil {
		return fmt.Errorf("system service not available")
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

func (e *TCPRTOExecutor) CanExecute(ctx context.Context) bool {
	elevated, err := e.elevationService.IsElevated(ctx)
	if err != nil || !elevated {
		return false
	}

	// Verificar se as configurações RTO já estão otimizadas
	registryPath := `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`
	
	currentRTT, err := e.registryService.ReadRegistryValue(ctx, registryPath, "TcpInitialRtt")
	if err != nil {
		return true // Pode ser configurado
	}

	// Se já está otimizado (valor 1000), não precisa executar
	if currentRTT == 1000 {
		return false
	}

	return true
}

func (e *TCPRTOExecutor) Revert(ctx context.Context, backupData entities.BackupData) (*entities.BoostRevertResult, error) {
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
			Message: "Administrator privileges required for TCP RTO revert",
			Error:   fmt.Errorf("elevation required"),
		}, fmt.Errorf("elevation required")
	}

	registryPath := `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`

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
		Message: "TCP retransmission interval reverted successfully",
		Error:   nil,
	}, nil
}

//go:build windows
// +build windows
package windows

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type WinServiceManagerService struct {
	logger *log.Logger
}

// NewWinServiceManagerService cria uma nova instância do service manager
func NewWinServiceManagerService(logger *log.Logger) *WinServiceManagerService {
	return &WinServiceManagerService{
		logger: logger,
	}
}

// StartService inicia um serviço do Windows
func (w *WinServiceManagerService) StartService(ctx context.Context, serviceName string) error {
	w.logger.Printf("Iniciando serviço: %s", serviceName)
	
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("erro ao conectar com o service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("erro ao abrir serviço '%s': %w", serviceName, err)
	}
	defer s.Close()

	// Verificar se já está rodando
	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("erro ao consultar status do serviço '%s': %w", serviceName, err)
	}

	if status.State == svc.Running {
		w.logger.Printf("Serviço '%s' já está rodando", serviceName)
		return nil
	}

	err = s.Start()
	if err != nil {
		return fmt.Errorf("erro ao iniciar serviço '%s': %w", serviceName, err)
	}

	// Aguardar o serviço iniciar
	err = w.waitForServiceState(s, svc.Running, 30*time.Second)
	if err != nil {
		return fmt.Errorf("timeout aguardando serviço '%s' iniciar: %w", serviceName, err)
	}

	w.logger.Printf("Serviço '%s' iniciado com sucesso", serviceName)
	return nil
}

// StopService para um serviço do Windows
func (w *WinServiceManagerService) StopService(ctx context.Context, serviceName string) error {
	w.logger.Printf("Parando serviço: %s", serviceName)

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("erro ao conectar com o service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("erro ao abrir serviço '%s': %w", serviceName, err)
	}
	defer s.Close()

	// Verificar se já está parado
	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("erro ao consultar status do serviço '%s': %w", serviceName, err)
	}

	if status.State == svc.Stopped {
		w.logger.Printf("Serviço '%s' já está parado", serviceName)
		return nil
	}

	if !status.Accepts.CanStop() {
		return fmt.Errorf("serviço '%s' não pode ser parado", serviceName)
	}

	_, err = s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("erro ao parar serviço '%s': %w", serviceName, err)
	}

	// Aguardar o serviço parar
	err = w.waitForServiceState(s, svc.Stopped, 30*time.Second)
	if err != nil {
		return fmt.Errorf("timeout aguardando serviço '%s' parar: %w", serviceName, err)
	}

	w.logger.Printf("Serviço '%s' parado com sucesso", serviceName)
	return nil
}

// RestartService reinicia um serviço do Windows
func (w *WinServiceManagerService) RestartService(ctx context.Context, serviceName string) error {
	w.logger.Printf("Reiniciando serviço: %s", serviceName)

	// Primeiro para o serviço
	err := w.StopService(ctx, serviceName)
	if err != nil {
		return fmt.Errorf("erro ao parar serviço para reiniciar: %w", err)
	}

	// Aguarda um pouco antes de iniciar
	time.Sleep(2 * time.Second)

	// Depois inicia o serviço
	err = w.StartService(ctx, serviceName)
	if err != nil {
		return fmt.Errorf("erro ao iniciar serviço após parar: %w", err)
	}

	w.logger.Printf("Serviço '%s' reiniciado com sucesso", serviceName)
	return nil
}

// SetServiceStartupType define o tipo de inicialização do serviço
func (w *WinServiceManagerService) SetServiceStartupType(ctx context.Context, serviceName string, startupType string) error {
	w.logger.Printf("Configurando tipo de inicialização do serviço '%s' para '%s'", serviceName, startupType)

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("erro ao conectar com o service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("erro ao abrir serviço '%s': %w", serviceName, err)
	}
	defer s.Close()

	var startType uint32
	switch strings.ToLower(startupType) {
	case "automatic":
		startType = mgr.StartAutomatic
	case "manual":
		startType = mgr.StartManual
	case "disabled":
		startType = mgr.StartDisabled
	case "delayed_automatic":
		startType = mgr.StartAutomatic // Será configurado como delayed depois
	default:
		return fmt.Errorf("tipo de inicialização inválido: %s", startupType)
	}

	config := mgr.Config{
		StartType: startType,
	}

	err = s.UpdateConfig(config)
	if err != nil {
		return fmt.Errorf("erro ao atualizar configuração do serviço '%s': %w", serviceName, err)
	}

	// Para delayed automatic, precisa de configuração adicional
	if strings.ToLower(startupType) == "delayed_automatic" {
		err = w.setDelayedAutoStart(s, true)
		if err != nil {
			w.logger.Printf("Aviso: Não foi possível configurar delayed auto start para '%s': %v", serviceName, err)
		}
	}

	w.logger.Printf("Tipo de inicialização do serviço '%s' configurado para '%s'", serviceName, startupType)
	return nil
}

// GetServiceStatus obtém o status detalhado de um serviço
func (w *WinServiceManagerService) GetServiceStatus(ctx context.Context, serviceName string) (*entities.WindowsService, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir serviço '%s': %w", serviceName, err)
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar status do serviço '%s': %w", serviceName, err)
	}

	config, err := s.Config()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter configuração do serviço '%s': %w", serviceName, err)
	}

	// Converter status para string
	statusStr := w.convertServiceState(status.State)
	startupTypeStr := w.convertStartupType(config.StartType)

	// Obter dependências
	dependencies, err := w.getServiceDependencies(s)
	if err != nil {
		w.logger.Printf("Aviso: Não foi possível obter dependências para '%s': %v", serviceName, err)
		dependencies = []string{}
	}

	service := &entities.WindowsService{
		ID:               fmt.Sprintf("win-svc-%s", serviceName),
		ServiceName:      serviceName,
		DisplayName:      config.DisplayName,
		Description:      config.Description,
		Status:           statusStr,
		StartupType:      startupTypeStr,
		IsEssential:      w.isEssentialService(serviceName),
		IsOptimized:      false, // Será definido por lógica de otimização
		CanBeStopped:     status.Accepts.CanStop(),
		CanBePaused:      status.Accepts.CanPauseContinue(),
		ProcessID:        int(status.ProcessId),
		MemoryUsage:      0, // Seria necessário consulta adicional
		Dependencies:     dependencies,
		DependentServices: []string{}, // Seria necessário consulta adicional
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return service, nil
}

// ListServices lista todos os serviços do Windows
func (w *WinServiceManagerService) ListServices(ctx context.Context) ([]*entities.WindowsService, error) {
	w.logger.Println("Listando todos os serviços do Windows")

	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o service manager: %w", err)
	}
	defer m.Disconnect()

	services, err := m.ListServices()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar serviços: %w", err)
	}

	var result []*entities.WindowsService
	for _, serviceName := range services {
		service, err := w.GetServiceStatus(ctx, serviceName)
		if err != nil {
			w.logger.Printf("Erro ao obter status do serviço '%s': %v", serviceName, err)
			continue
		}
		result = append(result, service)
	}

	w.logger.Printf("Listados %d serviços", len(result))
	return result, nil
}

// IsServiceRunning verifica se um serviço está rodando
func (w *WinServiceManagerService) IsServiceRunning(ctx context.Context, serviceName string) (bool, error) {
	service, err := w.GetServiceStatus(ctx, serviceName)
	if err != nil {
		return false, err
	}
	return service.Status == "running", nil
}

// DisableUnnecessaryServices desabilita serviços desnecessários para otimização
func (w *WinServiceManagerService) DisableUnnecessaryServices(ctx context.Context) error {
	w.logger.Println("Desabilitando serviços desnecessários")

	unnecessaryServices := []string{
		"Fax",
		"TabletInputService",
		"Spooler", // Apenas se não usa impressora
		"Themes",  // Para servidores
		"Windows Search",
		"Remote Registry",
		"Secondary Logon",
		"Routing and Remote Access",
	}

	for _, serviceName := range unnecessaryServices {
		err := w.SetServiceStartupType(ctx, serviceName, "disabled")
		if err != nil {
			w.logger.Printf("Aviso: Não foi possível desabilitar '%s': %v", serviceName, err)
			continue
		}

		// Tentar parar o serviço se estiver rodando
		isRunning, _ := w.IsServiceRunning(ctx, serviceName)
		if isRunning {
			err = w.StopService(ctx, serviceName)
			if err != nil {
				w.logger.Printf("Aviso: Não foi possível parar '%s': %v", serviceName, err)
			}
		}
	}

	w.logger.Println("Processo de desabilitação de serviços concluído")
	return nil
}

// EnableEssentialServices garante que serviços essenciais estejam habilitados
func (w *WinServiceManagerService) EnableEssentialServices(ctx context.Context) error {
	w.logger.Println("Habilitando serviços essenciais")

	essentialServices := map[string]string{
		"Winmgmt":      "automatic", // Windows Management Instrumentation
		"RpcSs":        "automatic", // Remote Procedure Call
		"DcomLaunch":   "automatic", // DCOM Server Process Launcher
		"PlugPlay":     "manual",    // Plug and Play
		"Power":        "automatic", // Power Service
		"Eventlog":     "automatic", // Windows Event Log
		"CryptSvc":     "automatic", // Cryptographic Services
		"Dhcp":         "automatic", // DHCP Client
		"Dnscache":     "automatic", // DNS Client
		"LanmanServer": "automatic", // Server
	}

	for serviceName, startupType := range essentialServices {
		err := w.SetServiceStartupType(ctx, serviceName, startupType)
		if err != nil {
			w.logger.Printf("Erro ao configurar serviço essencial '%s': %v", serviceName, err)
			continue
		}

		if startupType == "automatic" {
			err = w.StartService(ctx, serviceName)
			if err != nil {
				w.logger.Printf("Aviso: Não foi possível iniciar serviço essencial '%s': %v", serviceName, err)
			}
		}
	}

	w.logger.Println("Processo de habilitação de serviços essenciais concluído")
	return nil
}

// Métodos auxiliares

// waitForServiceState aguarda o serviço atingir um estado específico
func (w *WinServiceManagerService) waitForServiceState(s *mgr.Service, targetState svc.State, timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		status, err := s.Query()
		if err != nil {
			return err
		}
		if status.State == targetState {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout aguardando mudança de estado")
}

// convertServiceState converte o estado do serviço para string
func (w *WinServiceManagerService) convertServiceState(state svc.State) string {
	switch state {
	case svc.Running:
		return "running"
	case svc.Stopped:
		return "stopped"
	case svc.Paused:
		return "paused"
	case svc.StartPending:
		return "start_pending"
	case svc.StopPending:
		return "stop_pending"
	default:
		return "unknown"
	}
}

// convertStartupType converte o tipo de inicialização para string
func (w *WinServiceManagerService) convertStartupType(startType uint32) string {
	switch startType {
	case mgr.StartAutomatic:
		return "automatic"
	case mgr.StartManual:
		return "manual"
	case mgr.StartDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// isEssentialService determina se um serviço é essencial
func (w *WinServiceManagerService) isEssentialService(serviceName string) bool {
	essentialServices := map[string]bool{
		"Winmgmt":      true,
		"RpcSs":        true,
		"DcomLaunch":   true,
		"PlugPlay":     true,
		"Power":        true,
		"Eventlog":     true,
		"CryptSvc":     true,
		"Dhcp":         true,
		"Dnscache":     true,
		"LanmanServer": true,
		"Themes":       false,
		"Spooler":      false,
	}

	essential, exists := essentialServices[serviceName]
	return exists && essential
}

// getServiceDependencies obtém as dependências de um serviço
func (w *WinServiceManagerService) getServiceDependencies(s *mgr.Service) ([]string, error) {
	config, err := s.Config()
	if err != nil {
		return nil, err
	}
	return config.Dependencies, nil
}

// setDelayedAutoStart configura delayed automatic start (requer chamada Windows API)
func (w *WinServiceManagerService) setDelayedAutoStart(s *mgr.Service, delayed bool) error {
	// Esta implementação seria mais complexa e requereria chamadas diretas à Windows API
	// Por simplicidade, deixamos como stub
	w.logger.Printf("Configuração delayed auto start não implementada completamente")
	return nil
}

// Métodos adicionais para estatísticas e relatórios

// GetServiceStatistics retorna estatísticas dos serviços
func (w *WinServiceManagerService) GetServiceStatistics(ctx context.Context) (map[string]interface{}, error) {
	services, err := w.ListServices(ctx)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_services": len(services),
		"running":        0,
		"stopped":        0,
		"essential":      0,
		"optimized":      0,
	}

	for _, service := range services {
		if service.Status == "running" {
			stats["running"] = stats["running"].(int) + 1
		}
		if service.Status == "stopped" {
			stats["stopped"] = stats["stopped"].(int) + 1
		}
		if service.IsEssential {
			stats["essential"] = stats["essential"].(int) + 1
		}
		if service.IsOptimized {
			stats["optimized"] = stats["optimized"].(int) + 1
		}
	}

	return stats, nil
}
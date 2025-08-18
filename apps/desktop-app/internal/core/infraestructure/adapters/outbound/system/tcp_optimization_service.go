//go:build windows
// +build windows
package system

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

// TCPOptimizationServiceImpl implementa TCPOptimizationService para Windows
type TCPOptimizationServiceImpl struct {
	// Registry paths para configurações TCP
	tcpParametersPath string
	// Cache das configurações originais para restauração
	originalConfig *entities.TCPConfiguration
}

// NewTCPOptimizationService cria uma nova instância do serviço
func NewTCPOptimizationService() *TCPOptimizationServiceImpl {
	return &TCPOptimizationServiceImpl{
		tcpParametersPath: `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
	}
}

// OptimizeTCPForGaming otimiza TCP para jogos (baixa latência)
func (s *TCPOptimizationServiceImpl) OptimizeTCPForGaming(ctx context.Context) error {
	// Salvar configuração atual antes de modificar
	if s.originalConfig == nil {
		config, err := s.GetTCPConfiguration(ctx)
		if err != nil {
			return fmt.Errorf("failed to backup current configuration: %w", err)
		}
		s.originalConfig = config
	}

	optimizations := map[string]interface{}{
		// Disable Nagle Algorithm para reduzir latência
		"TcpNoDelay":           uint32(1),
		// Otimizar TCP Window Size
		"TcpWindowSize":        uint32(65536), // 64KB
		// TCP Chimney Offload (se suportado)
		"EnableTCPChimney":     uint32(1),
		// Disable TCP Auto-Tuning Level (para controle manual)
		"TcpAutoTuningLevel":   uint32(0),
		// RSS (Receive Side Scaling)
		"EnableRSS":            uint32(1),
		// Interrupt Moderation
		"InterruptModeration":  uint32(1),
	}

	return s.applyTCPOptimizations(ctx, optimizations)
}

// OptimizeTCPForStreaming otimiza TCP para streaming (throughput)
func (s *TCPOptimizationServiceImpl) OptimizeTCPForStreaming(ctx context.Context) error {
	// Salvar configuração atual
	if s.originalConfig == nil {
		config, err := s.GetTCPConfiguration(ctx)
		if err != nil {
			return fmt.Errorf("failed to backup current configuration: %w", err)
		}
		s.originalConfig = config
	}

	optimizations := map[string]interface{}{
		// Larger Window Size para throughput
		"TcpWindowSize":        uint32(131072), // 128KB
		// Enable TCP Auto-Tuning
		"TcpAutoTuningLevel":   uint32(2), // Normal
		// Congestion Control otimizado
		"CongestionProvider":   "cubic",
		// Buffer sizes maiores
		"DefaultRcvWindow":     uint32(65536),
		"DefaultSendWindow":    uint32(65536),
		// RSS habilitado
		"EnableRSS":            uint32(1),
	}

	return s.applyTCPOptimizations(ctx, optimizations)
}

// RestoreDefaultTCPSettings restaura configurações padrão
func (s *TCPOptimizationServiceImpl) RestoreDefaultTCPSettings(ctx context.Context) error {
	if s.originalConfig == nil {
		return fmt.Errorf("no backup configuration available")
	}

	// Converter configuração de volta para registry values
	defaultSettings := map[string]interface{}{
		"TcpNoDelay":          uint32(0),
		"TcpWindowSize":       uint32(s.originalConfig.WindowSize),
		"TcpAutoTuningLevel":  uint32(2), // Normal (padrão Windows)
		"EnableTCPChimney":    uint32(0),
		"EnableRSS":           uint32(1),
	}

	err := s.applyTCPOptimizations(ctx, defaultSettings)
	if err != nil {
		return err
	}

	// Reset congestion control to default
	return s.SetTCPCongestionControl(ctx, "cubic")
}

// SetTCPWindowSize define o tamanho da janela TCP
func (s *TCPOptimizationServiceImpl) SetTCPWindowSize(ctx context.Context, size int) error {
	if size < 8192 || size > 16777216 { // 8KB to 16MB
		return fmt.Errorf("invalid window size: %d (must be between 8192 and 16777216)", size)
	}

	return s.setRegistryValue("TcpWindowSize", uint32(size))
}

// SetTCPCongestionControl define o algoritmo de controle de congestionamento
func (s *TCPOptimizationServiceImpl) SetTCPCongestionControl(ctx context.Context, algorithm string) error {
	validAlgorithms := []string{"cubic", "bbr", "reno", "newreno"}
	
	algorithm = strings.ToLower(algorithm)
	isValid := false
	for _, valid := range validAlgorithms {
		if algorithm == valid {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return fmt.Errorf("invalid congestion control algorithm: %s", algorithm)
	}

	// Usar netsh para configurar congestion control
	cmd := exec.CommandContext(ctx, "netsh", "int", "tcp", "set", "global", 
		fmt.Sprintf("congestionprovider=%s", algorithm))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set congestion control: %s, output: %s", err, string(output))
	}

	return nil
}

// DisableNagleAlgorithm desabilita o algoritmo de Nagle
func (s *TCPOptimizationServiceImpl) DisableNagleAlgorithm(ctx context.Context) error {
	return s.setRegistryValue("TcpNoDelay", uint32(1))
}

// GetTCPConfiguration retorna a configuração atual do TCP
func (s *TCPOptimizationServiceImpl) GetTCPConfiguration(ctx context.Context) (*entities.TCPConfiguration, error) {
	config := &entities.TCPConfiguration{}

	// Ler valores do registry
	windowSize, err := s.getRegistryUint32Value("TcpWindowSize")
	if err == nil {
		config.WindowSize = int(windowSize)
	} else {
		config.WindowSize = 65536 // Default
	}

	nagleDisabled, err := s.getRegistryUint32Value("TcpNoDelay")
	if err == nil {
		config.NagleDisabled = nagleDisabled == 1
	}

	autoTuning, err := s.getRegistryUint32Value("TcpAutoTuningLevel")
	if err == nil {
		config.AutoTuning = autoTuning != 0
	}

	// Obter congestion control via netsh
	congestionControl, err := s.getCurrentCongestionControl(ctx)
	if err == nil {
		config.CongestionControl = congestionControl
	}

	// Detectar se está otimizado baseado nas configurações
	config.Optimized = s.isConfigurationOptimized(config)

	return config, nil
}

// IsTCPOptimized verifica se o TCP está otimizado
func (s *TCPOptimizationServiceImpl) IsTCPOptimized(ctx context.Context) (bool, error) {
	config, err := s.GetTCPConfiguration(ctx)
	if err != nil {
		return false, err
	}

	return config.Optimized, nil
}

// GetNetworkLatency mede a latência para um alvo específico
func (s *TCPOptimizationServiceImpl) GetNetworkLatency(ctx context.Context, target string) (int, error) {
	// Usar ping nativo do Windows para medir latência
	cmd := exec.CommandContext(ctx, "ping", "-n", "4", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to ping %s: %w", target, err)
	}

	// Extrair latência média do output do ping
	return s.parsePingLatency(string(output))
}

// Métodos auxiliares privados

func (s *TCPOptimizationServiceImpl) applyTCPOptimizations(ctx context.Context, optimizations map[string]interface{}) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, s.tcpParametersPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	for name, value := range optimizations {
		switch v := value.(type) {
		case uint32:
			err = key.SetDWordValue(name, v)
		case string:
			if name == "CongestionProvider" {
				// Usar netsh para congestion control
				err = s.SetTCPCongestionControl(ctx, v)
				continue
			}
			err = key.SetStringValue(name, v)
		default:
			err = fmt.Errorf("unsupported value type for %s", name)
		}

		if err != nil {
			return fmt.Errorf("failed to set registry value %s: %w", name, err)
		}
	}

	return nil
}

func (s *TCPOptimizationServiceImpl) setRegistryValue(name string, value uint32) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, s.tcpParametersPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	return key.SetDWordValue(name, value)
}

func (s *TCPOptimizationServiceImpl) getRegistryUint32Value(name string) (uint32, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, s.tcpParametersPath, registry.QUERY_VALUE)
	if err != nil {
		return 0, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue(name)
	if err != nil {
		return 0, fmt.Errorf("failed to read registry value %s: %w", name, err)
	}

	return uint32(value), nil
}

func (s *TCPOptimizationServiceImpl) getCurrentCongestionControl(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "netsh", "int", "tcp", "show", "global")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get TCP global settings: %w", err)
	}

	// Extrair congestion provider do output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "congestion provider") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "cubic", nil // Default assumption
}

func (s *TCPOptimizationServiceImpl) isConfigurationOptimized(config *entities.TCPConfiguration) bool {
	// Considera otimizado se:
	// - Window size > 32KB
	// - Nagle desabilitado OU congestion control não é padrão
	// - Configurações específicas para gaming ou streaming
	
	return config.WindowSize > 32768 && 
		   (config.NagleDisabled || 
		    (config.CongestionControl != "cubic" && config.CongestionControl != ""))
}

func (s *TCPOptimizationServiceImpl) parsePingLatency(output string) (int, error) {
	// Regex para extrair tempo médio do ping do Windows
	// Exemplo: "Average = 45ms"
	re := regexp.MustCompile(`Average = (\d+)ms`)
	matches := re.FindStringSubmatch(output)
	
	if len(matches) >= 2 {
		latency, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("failed to parse latency: %w", err)
		}
		return latency, nil
	}

	// Fallback: tentar extrair de uma única linha de ping
	re = regexp.MustCompile(`time=(\d+)ms`)
	matches = re.FindStringSubmatch(output)
	
	if len(matches) >= 2 {
		latency, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("failed to parse latency: %w", err)
		}
		return latency, nil
	}

	return 0, fmt.Errorf("could not parse ping latency from output")
}

//go:build windows
// +build windows
package system

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// Windows Power Management constants
const (
	POWER_SETTING_GUID_GPU_PREFERENCE = "{5d76a2ca-e8c0-402f-a133-2158492d58ad}"
	POWER_SCHEME_GUID_HIGH_PERFORMANCE = "{8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c}"
	POWER_SCHEME_GUID_BALANCED         = "{381b4222-f694-41f0-9685-ff5bb260df2e}"
	POWER_SCHEME_GUID_POWER_SAVER      = "{a1841308-3541-4fab-bc81-f71556f20b4a}"
)

// GPU power modes
const (
	PowerModeMaxPerformance = "max_performance"
	PowerModeBalanced       = "balanced"
	PowerModePowerSaver     = "power_saver"
	PowerModeOptimal        = "optimal"
)

// Registry paths for GPU settings
const (
	NVIDIARegistryPath    = `SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}`
	AMDRegistryPath       = `SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}`
	GraphicsDriversPath   = `SOFTWARE\Microsoft\DirectX\UserGpuPreferences`
	PowerSettingsPath     = `SYSTEM\CurrentControlSet\Control\Power\User\PowerSchemes`
)

// GPUOptimizationService implementation for Windows
type gpuOptimizationService struct {
	mu                sync.RWMutex
	currentConfig     *entities.GPUConfiguration
	originalConfig    *entities.GPUConfiguration
	isOptimized       bool
	gpuRepository     GPUInfoRepository
	powerAPI          *PowerAPI
}

// PowerAPI encapsula chamadas para APIs de energia do Windows
type PowerAPI struct {
	powrprof *windows.LazyDLL
}

// NewGPUOptimizationService creates a new Windows GPU optimization service
func NewGPUOptimizationService(gpuRepo GPUInfoRepository) *gpuOptimizationService {
	return &gpuOptimizationService{
		gpuRepository: gpuRepo,
		powerAPI: &PowerAPI{
			powrprof: windows.NewLazySystemDLL("powrprof.dll"),
		},
	}
}

// OptimizeGPUForGaming otimiza a GPU para jogos
func (s *gpuOptimizationService) OptimizeGPUForGaming(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Salvar configuração original se ainda não foi salva
	if s.originalConfig == nil {
		if err := s.saveOriginalConfiguration(ctx); err != nil {
			return fmt.Errorf("failed to save original configuration: %w", err)
		}
	}

	// Aplicar configurações de gaming
	optimizations := []func(context.Context) error{
		s.setHighPerformancePowerPlan,
		s.setMaxGPUPerformanceMode,
		s.disableGPUPowerSaving,
		s.optimizeGPUScheduling,
		s.setGamingPreferences,
	}

	for _, opt := range optimizations {
		if err := opt(ctx); err != nil {
			return fmt.Errorf("gaming optimization failed: %w", err)
		}
	}

	s.isOptimized = true
	s.currentConfig = &entities.GPUConfiguration{
		PowerMode:        PowerModeMaxPerformance,
		OptimizationMode: "gaming",
		IsOverclocked:    false,
		UpdatedAt:        time.Now(),
	}

	return nil
}

// OptimizeGPUForWorkstation otimiza a GPU para workstation
func (s *gpuOptimizationService) OptimizeGPUForWorkstation(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Salvar configuração original se ainda não foi salva
	if s.originalConfig == nil {
		if err := s.saveOriginalConfiguration(ctx); err != nil {
			return fmt.Errorf("failed to save original configuration: %w", err)
		}
	}

	// Aplicar configurações de workstation
	optimizations := []func(context.Context) error{
		s.setBalancedPowerPlan,
		s.setOptimalGPUPerformanceMode,
		s.enableGPUCompute,
		s.optimizeMemoryAllocation,
		s.setWorkstationPreferences,
	}

	for _, opt := range optimizations {
		if err := opt(ctx); err != nil {
			return fmt.Errorf("workstation optimization failed: %w", err)
		}
	}

	s.isOptimized = true
	s.currentConfig = &entities.GPUConfiguration{
		PowerMode:        PowerModeOptimal,
		OptimizationMode: "workstation",
		IsOverclocked:    false,
		UpdatedAt:        time.Now(),
	}

	return nil
}

// RestoreGPUDefaults restaura as configurações padrão da GPU
func (s *gpuOptimizationService) RestoreGPUDefaults(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.originalConfig == nil {
		return fmt.Errorf("no original configuration to restore")
	}

	// Restaurar configurações originais
	restorations := []func(context.Context) error{
		s.restorePowerPlan,
		s.restoreGPUSettings,
		s.restoreRegistrySettings,
		s.clearOptimizationFlags,
	}

	for _, restore := range restorations {
		if err := restore(ctx); err != nil {
			return fmt.Errorf("restoration failed: %w", err)
		}
	}

	s.isOptimized = false
	s.currentConfig = s.originalConfig

	return nil
}

// SetGPUPowerMode define o modo de energia da GPU
func (s *gpuOptimizationService) SetGPUPowerMode(ctx context.Context, mode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var powerGUID string
	switch strings.ToLower(mode) {
	case PowerModeMaxPerformance:
		powerGUID = POWER_SCHEME_GUID_HIGH_PERFORMANCE
	case PowerModeBalanced:
		powerGUID = POWER_SCHEME_GUID_BALANCED
	case PowerModePowerSaver:
		powerGUID = POWER_SCHEME_GUID_POWER_SAVER
	default:
		return fmt.Errorf("invalid power mode: %s", mode)
	}

	if err := s.setPowerScheme(powerGUID); err != nil {
		return fmt.Errorf("failed to set power mode: %w", err)
	}

	if s.currentConfig != nil {
		s.currentConfig.PowerMode = mode
		s.currentConfig.UpdatedAt = time.Now()
	}

	return nil
}

// SetGPUMemoryClock define a velocidade do clock de memória
func (s *gpuOptimizationService) SetGPUMemoryClock(ctx context.Context, clockSpeed int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if clockSpeed < 0 {
		return fmt.Errorf("invalid memory clock speed: %d", clockSpeed)
	}

	// Tentar diferentes métodos baseados no vendor da GPU
	gpus, err := s.gpuRepository.GetGPUInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get GPU info: %w", err)
	}

	for _, gpu := range gpus {
		switch strings.ToLower(gpu.Vendor) {
		case "nvidia":
			if err := s.setNVIDIAMemoryClock(gpu.DeviceID, clockSpeed); err != nil {
				return fmt.Errorf("failed to set NVIDIA memory clock: %w", err)
			}
		case "amd":
			if err := s.setAMDMemoryClock(gpu.DeviceID, clockSpeed); err != nil {
				return fmt.Errorf("failed to set AMD memory clock: %w", err)
			}
		default:
			return fmt.Errorf("memory clock adjustment not supported for vendor: %s", gpu.Vendor)
		}
	}

	if s.currentConfig != nil {
		s.currentConfig.MemoryClock = clockSpeed
		s.currentConfig.IsOverclocked = true
		s.currentConfig.UpdatedAt = time.Now()
	}

	return nil
}

// SetGPUCoreClock define a velocidade do clock do core
func (s *gpuOptimizationService) SetGPUCoreClock(ctx context.Context, clockSpeed int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if clockSpeed < 0 {
		return fmt.Errorf("invalid core clock speed: %d", clockSpeed)
	}

	gpus, err := s.gpuRepository.GetGPUInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get GPU info: %w", err)
	}

	for _, gpu := range gpus {
		switch strings.ToLower(gpu.Vendor) {
		case "nvidia":
			if err := s.setNVIDIACoreClock(gpu.DeviceID, clockSpeed); err != nil {
				return fmt.Errorf("failed to set NVIDIA core clock: %w", err)
			}
		case "amd":
			if err := s.setAMDCoreClock(gpu.DeviceID, clockSpeed); err != nil {
				return fmt.Errorf("failed to set AMD core clock: %w", err)
			}
		default:
			return fmt.Errorf("core clock adjustment not supported for vendor: %s", gpu.Vendor)
		}
	}

	if s.currentConfig != nil {
		s.currentConfig.CoreClock = clockSpeed
		s.currentConfig.IsOverclocked = true
		s.currentConfig.UpdatedAt = time.Now()
	}

	return nil
}

// GetGPUConfiguration retorna a configuração atual da GPU
func (s *gpuOptimizationService) GetGPUConfiguration(ctx context.Context) (*entities.GPUConfiguration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.currentConfig == nil {
		// Detectar configuração atual
		config, err := s.detectCurrentConfiguration(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to detect current configuration: %w", err)
		}
		return config, nil
	}

	// Retornar cópia da configuração atual
	return &entities.GPUConfiguration{
		PowerMode:        s.currentConfig.PowerMode,
		OptimizationMode: s.currentConfig.OptimizationMode,
		CoreClock:        s.currentConfig.CoreClock,
		MemoryClock:      s.currentConfig.MemoryClock,
		IsOverclocked:    s.currentConfig.IsOverclocked,
		CreatedAt:        s.currentConfig.CreatedAt,
		UpdatedAt:        s.currentConfig.UpdatedAt,
	}, nil
}

// IsGPUOptimized verifica se a GPU está otimizada
func (s *gpuOptimizationService) IsGPUOptimized(ctx context.Context) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.isOptimized, nil
}

// Implementações dos métodos de otimização

func (s *gpuOptimizationService) saveOriginalConfiguration(ctx context.Context) error {
	config, err := s.detectCurrentConfiguration(ctx)
	if err != nil {
		return err
	}
	s.originalConfig = config
	return nil
}

func (s *gpuOptimizationService) setHighPerformancePowerPlan(ctx context.Context) error {
	return s.setPowerScheme(POWER_SCHEME_GUID_HIGH_PERFORMANCE)
}

func (s *gpuOptimizationService) setBalancedPowerPlan(ctx context.Context) error {
	return s.setPowerScheme(POWER_SCHEME_GUID_BALANCED)
}

func (s *gpuOptimizationService) setPowerScheme(schemeGUID string) error {
	// Usar powercfg via registry para definir o esquema de energia
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, PowerSettingsPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open power settings key: %w", err)
	}
	defer key.Close()

	return key.SetStringValue("ActivePowerScheme", schemeGUID)
}

func (s *gpuOptimizationService) setMaxGPUPerformanceMode(ctx context.Context) error {
	return s.setGPUPreference("GpuPreference", "2") // 2 = High Performance
}

func (s *gpuOptimizationService) setOptimalGPUPerformanceMode(ctx context.Context) error {
	return s.setGPUPreference("GpuPreference", "1") // 1 = Auto Select
}

func (s *gpuOptimizationService) setGPUPreference(setting, value string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, GraphicsDriversPath, registry.CREATE_SUB_KEY|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open graphics drivers key: %w", err)
	}
	defer key.Close()

	return key.SetStringValue(setting, value)
}

func (s *gpuOptimizationService) disableGPUPowerSaving(ctx context.Context) error {
	// Desabilitar power saving features via WMI
	return s.executeWMICommand(`
		$gpu = Get-WmiObject -Class Win32_VideoController | Where-Object {$_.VideoProcessor -ne $null}
		foreach($g in $gpu) {
			$g.SetPowerState(1)  # Full power
		}
	`)
}

func (s *gpuOptimizationService) enableGPUCompute(ctx context.Context) error {
	// Habilitar compute shaders e otimizações para workstation
	return s.setRegistryValue(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\DirectX`,
		"EnableComputeShaders",
		uint32(1),
	)
}

func (s *gpuOptimizationService) optimizeGPUScheduling(ctx context.Context) error {
	// Otimizar GPU scheduling (Windows 10/11)
	return s.setRegistryValue(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\GraphicsDrivers`,
		"HwSchMode",
		uint32(2), // Hardware-accelerated GPU scheduling
	)
}

func (s *gpuOptimizationService) optimizeMemoryAllocation(ctx context.Context) error {
	// Otimizar alocação de memória de vídeo
	return s.setRegistryValue(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\GraphicsDrivers`,
		"DedicatedSegmentSize",
		uint32(0), // Let system manage
	)
}

func (s *gpuOptimizationService) setGamingPreferences(ctx context.Context) error {
	preferences := map[string]interface{}{
		"GameMode":           uint32(1),
		"AutoGameMode":       uint32(1),
		"GameDVR_Enabled":    uint32(0), // Disable GameDVR for performance
	}

	for key, value := range preferences {
		if err := s.setRegistryValue(registry.CURRENT_USER, `SOFTWARE\Microsoft\GameBar`, key, value); err != nil {
			return err
		}
	}

	return nil
}

func (s *gpuOptimizationService) setWorkstationPreferences(ctx context.Context) error {
	// Configurações otimizadas para workstation
	preferences := map[string]interface{}{
		"ApplicationProfile":     "Workstation",
		"PreferMaxPerformance":   uint32(1),
		"EnableGPUCompute":       uint32(1),
	}

	for key, value := range preferences {
		if err := s.setRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\DirectX`, key, value); err != nil {
			return err
		}
	}

	return nil
}

// Métodos de restauração

func (s *gpuOptimizationService) restorePowerPlan(ctx context.Context) error {
	if s.originalConfig != nil && s.originalConfig.PowerMode != "" {
		return s.SetGPUPowerMode(ctx, s.originalConfig.PowerMode)
	}
	return s.setPowerScheme(POWER_SCHEME_GUID_BALANCED) // Default
}

func (s *gpuOptimizationService) restoreGPUSettings(ctx context.Context) error {
	// Restaurar configurações de GPU para padrão
	return s.setGPUPreference("GpuPreference", "0") // 0 = Auto
}

func (s *gpuOptimizationService) restoreRegistrySettings(ctx context.Context) error {
	// Restaurar configurações críticas do registro
	settings := []struct {
		hkey   registry.Key
		path   string
		name   string
		value  interface{}
	}{
		{registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\GraphicsDrivers`, "HwSchMode", uint32(1)},
		{registry.CURRENT_USER, `SOFTWARE\Microsoft\GameBar`, "GameMode", uint32(0)},
	}

	for _, setting := range settings {
		s.setRegistryValue(setting.hkey, setting.path, setting.name, setting.value)
	}

	return nil
}

func (s *gpuOptimizationService) clearOptimizationFlags(ctx context.Context) error {
	s.isOptimized = false
	return nil
}

// Métodos específicos de vendor

func (s *gpuOptimizationService) setNVIDIAMemoryClock(deviceID string, clockSpeed int) error {
	// Implementação simplificada - em produção usar NVAPI
	return s.setRegistryValue(
		registry.LOCAL_MACHINE,
		fmt.Sprintf(`%s\0000`, NVIDIARegistryPath),
		"MemoryClockOffset",
		uint32(clockSpeed),
	)
}

func (s *gpuOptimizationService) setNVIDIACoreClock(deviceID string, clockSpeed int) error {
	return s.setRegistryValue(
		registry.LOCAL_MACHINE,
		fmt.Sprintf(`%s\0000`, NVIDIARegistryPath),
		"CoreClockOffset",
		uint32(clockSpeed),
	)
}

func (s *gpuOptimizationService) setAMDMemoryClock(deviceID string, clockSpeed int) error {
	return s.setRegistryValue(
		registry.LOCAL_MACHINE,
		fmt.Sprintf(`%s\0000`, AMDRegistryPath),
		"PP_MemoryClockOverride",
		uint32(clockSpeed),
	)
}

func (s *gpuOptimizationService) setAMDCoreClock(deviceID string, clockSpeed int) error {
	return s.setRegistryValue(
		registry.LOCAL_MACHINE,
		fmt.Sprintf(`%s\0000`, AMDRegistryPath),
		"PP_GPUClockOverride",
		uint32(clockSpeed),
	)
}

// Métodos auxiliares

func (s *gpuOptimizationService) detectCurrentConfiguration(ctx context.Context) (*entities.GPUConfiguration, error) {
	now := time.Now()
	
	// Detectar modo de energia atual
	powerMode, err := s.getCurrentPowerMode()
	if err != nil {
		powerMode = PowerModeBalanced // Default
	}

	return &entities.GPUConfiguration{
		PowerMode:        powerMode,
		OptimizationMode: "default",
		CoreClock:        0,
		MemoryClock:      0,
		IsOverclocked:    false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func (s *gpuOptimizationService) getCurrentPowerMode() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, PowerSettingsPath, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	scheme, _, err := key.GetStringValue("ActivePowerScheme")
	if err != nil {
		return "", err
	}

	switch scheme {
	case POWER_SCHEME_GUID_HIGH_PERFORMANCE:
		return PowerModeMaxPerformance, nil
	case POWER_SCHEME_GUID_POWER_SAVER:
		return PowerModePowerSaver, nil
	default:
		return PowerModeBalanced, nil
	}
}

func (s *gpuOptimizationService) setRegistryValue(hkey registry.Key, path, name string, value interface{}) error {
	key, err := registry.OpenKey(hkey, path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key %s: %w", path, err)
	}
	defer key.Close()

	switch v := value.(type) {
	case string:
		return key.SetStringValue(name, v)
	case uint32:
		return key.SetDWordValue(name, v)
	case uint64:
		return key.SetQWordValue(name, v)
	default:
		return fmt.Errorf("unsupported registry value type: %T", value)
	}
}

func (s *gpuOptimizationService) executeWMICommand(command string) error {
	// Implementação simplificada usando WMI via COM
	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer unknown.Release()

	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer shell.Release()

	_, err = oleutil.CallMethod(shell, "Run", 
		fmt.Sprintf(`powershell.exe -Command "%s"`, command), 0, true)
	
	return err
}
//go:build windows
// +build windows
package system

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

// WindowsPowerManagementService implementa o serviço de gerenciamento de energia para Windows
type WindowsPowerManagementService struct {
	ole     *ole.IUnknown
	wmi     *ole.IDispatch
	locator *ole.IDispatch
}

// powerSchemeGUIDs contém os GUIDs dos perfis de energia padrão do Windows
var powerSchemeGUIDs = map[string]string{
	"SCHEME_BALANCED":     "381b4222-f694-41f0-9685-ff5bb260df2e",
	"SCHEME_MAX_POWER":    "8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c",
	"SCHEME_POWER_SAVER":  "a1841308-3541-4fab-bc81-f71556f20b4a",
}

// NewWindowsPowerManagementService cria uma nova instância do serviço de energia
func NewWindowsPowerManagementService() (*WindowsPowerManagementService, error) {
	err := ole.CoInitialize(0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize COM: %w", err)
	}

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		ole.CoUninitialize()
		return nil, fmt.Errorf("failed to create WMI locator: %w", err)
	}

	locator, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		unknown.Release()
		ole.CoUninitialize()
		return nil, fmt.Errorf("failed to query WMI locator interface: %w", err)
	}

	// Conectar ao namespace WMI root\cimv2\power
	serviceRaw, err := oleutil.CallMethod(locator, "ConnectServer", nil, "root\\cimv2\\power")
	if err != nil {
		locator.Release()
		unknown.Release()
		ole.CoUninitialize()
		return nil, fmt.Errorf("failed to connect to WMI power namespace: %w", err)
	}

	service := serviceRaw.ToIDispatch()

	return &WindowsPowerManagementService{
		ole:     unknown,
		wmi:     service,
		locator: locator,
	}, nil
}

// Close libera os recursos COM
func (w *WindowsPowerManagementService) Close() {
	if w.wmi != nil {
		w.wmi.Release()
	}
	if w.locator != nil {
		w.locator.Release()
	}
	if w.ole != nil {
		w.ole.Release()
	}
	ole.CoUninitialize()
}

// SetPowerProfile define o perfil de energia ativo
func (w *WindowsPowerManagementService) SetPowerProfile(ctx context.Context, profileName string) error {
	// Primeiro, tenta mapear para GUID conhecido
	var profileGUID string
	if guid, exists := powerSchemeGUIDs[strings.ToUpper(profileName)]; exists {
		profileGUID = guid
	} else {
		// Se não for um perfil padrão, busca pelo nome
		profiles, err := w.GetAvailablePowerProfiles(ctx)
		if err != nil {
			return fmt.Errorf("failed to get available profiles: %w", err)
		}

		for _, profile := range profiles {
			if strings.EqualFold(profile.Name, profileName) {
				profileGUID = profile.ID
				break
			}
		}

		if profileGUID == "" {
			return fmt.Errorf("power profile '%s' not found", profileName)
		}
	}

	// Usar powercfg para definir o perfil ativo
	cmd := exec.CommandContext(ctx, "powercfg", "/setactive", profileGUID)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set power profile: %w", err)
	}

	return nil
}

// GetActivePowerProfile retorna o perfil de energia ativo
func (w *WindowsPowerManagementService) GetActivePowerProfile(ctx context.Context) (*entities.PowerProfile, error) {
	cmd := exec.CommandContext(ctx, "powercfg", "/getactivescheme")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get active power scheme: %w", err)
	}

	// Parse da saída: "Power Scheme GUID: <guid> (<name>)"
	line := strings.TrimSpace(string(output))
	parts := strings.Split(line, " ")
	
	var guid, name string
	for _, part := range parts {
		if strings.Contains(part, "-") && len(part) == 36 {
			guid = part
			// Nome está entre parênteses no final
			nameStart := strings.Index(line, "(")
			nameEnd := strings.LastIndex(line, ")")
			if nameStart != -1 && nameEnd != -1 && nameEnd > nameStart {
				name = line[nameStart+1:nameEnd]
			}
			break
		}
	}

	if guid == "" {
		return nil, fmt.Errorf("could not parse active power scheme")
	}

	return &entities.PowerProfile{
		ID:          guid,
		Name:        name,
		Description: fmt.Sprintf("Active power profile: %s", name),
		IsActive:    true,
	}, nil
}

// GetAvailablePowerProfiles retorna todos os perfis de energia disponíveis
func (w *WindowsPowerManagementService) GetAvailablePowerProfiles(ctx context.Context) ([]*entities.PowerProfile, error) {
	cmd := exec.CommandContext(ctx, "powercfg", "/list")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list power schemes: %w", err)
	}

	var profiles []*entities.PowerProfile
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Power Scheme GUID:") {
			// Parse: "Power Scheme GUID: <guid> (<name>) *"
			parts := strings.Split(line, " ")
			var guid, name string
			isActive := strings.Contains(line, " *")

			for _, part := range parts {
				if strings.Contains(part, "-") && len(part) == 36 {
					guid = part
					break
				}
			}

			nameStart := strings.Index(line, "(")
			nameEnd := strings.Index(line, ")")
			if nameStart != -1 && nameEnd != -1 && nameEnd > nameStart {
				name = line[nameStart+1:nameEnd]
				if isActive {
					name = strings.TrimSuffix(name, " *")
				}
			}

			if guid != "" && name != "" {
				profiles = append(profiles, &entities.PowerProfile{
					ID:          guid,
					Name:        name,
					Description: fmt.Sprintf("Power profile: %s", name),
					IsActive:    isActive,
				})
			}
		}
	}

	return profiles, nil
}

// DisablePowerThrottling desabilita o throttling de energia
func (w *WindowsPowerManagementService) DisablePowerThrottling(ctx context.Context) error {
	// Desabilita o throttling de energia via registro do Windows
	commands := [][]string{
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_PROCESSOR", "PERFINCTHRESHOLD", "0"},
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_PROCESSOR", "PERFDECTHRESHOLD", "0"},
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_PROCESSOR", "PERFINCTIME", "1"},
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_PROCESSOR", "PERFDECTIME", "1"},
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_PROCESSOR", "PROCTHROTTLEMAX", "100"},
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_PROCESSOR", "PROCTHROTTLEMIN", "100"},
		{"powercfg", "/setactive", "SCHEME_CURRENT"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to execute command %v: %w", cmdArgs, err)
		}
	}

	return nil
}

// SetCPUPowerPolicy define a política de energia da CPU
func (w *WindowsPowerManagementService) SetCPUPowerPolicy(ctx context.Context, policy string) error {
	var policyValue string

	switch strings.ToLower(policy) {
	case "performance", "high":
		policyValue = "100"
	case "balanced":
		policyValue = "50"
	case "power_saver", "low":
		policyValue = "5"
	default:
		// Tenta interpretar como valor numérico
		if val, err := strconv.Atoi(policy); err == nil && val >= 0 && val <= 100 {
			policyValue = policy
		} else {
			return fmt.Errorf("invalid CPU power policy: %s. Use 'performance', 'balanced', 'power_saver', or a value 0-100", policy)
		}
	}

	commands := [][]string{
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_PROCESSOR", "PROCTHROTTLEMAX", policyValue},
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_PROCESSOR", "PROCTHROTTLEMIN", policyValue},
		{"powercfg", "/setactive", "SCHEME_CURRENT"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set CPU power policy: %w", err)
		}
	}

	return nil
}

// SetUSBPowerSettings configura as configurações de energia USB
func (w *WindowsPowerManagementService) SetUSBPowerSettings(ctx context.Context, enabled bool) error {
	value := "1"
	if !enabled {
		value = "0"
	}

	commands := [][]string{
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", "SUB_USB", "USBSELECTIVESUSPEND", value},
		{"powercfg", "/setdcvalueindex", "SCHEME_CURRENT", "SUB_USB", "USBSELECTIVESUSPEND", value},
		{"powercfg", "/setactive", "SCHEME_CURRENT"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set USB power settings: %w", err)
		}
	}

	return nil
}

// CreateCustomPowerProfile cria um perfil de energia personalizado
func (w *WindowsPowerManagementService) CreateCustomPowerProfile(ctx context.Context, profile *entities.PowerProfile) error {
	// Duplica um perfil existente como base
	baseScheme := "SCHEME_BALANCED"
	if profile.BaseProfile != "" {
		if guid, exists := powerSchemeGUIDs[strings.ToUpper(profile.BaseProfile)]; exists {
			baseScheme = guid
		} else {
			baseScheme = profile.BaseProfile
		}
	}

	// Cria o perfil duplicando um existente
	cmd := exec.CommandContext(ctx, "powercfg", "/duplicatescheme", baseScheme, profile.ID)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create custom power profile: %w", err)
	}

	// Define o nome do perfil
	cmd = exec.CommandContext(ctx, "powercfg", "/changename", profile.ID, profile.Name, profile.Description)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		// Tenta remover o perfil se não conseguiu definir o nome
		w.DeletePowerProfile(ctx, profile.ID)
		return fmt.Errorf("failed to set custom power profile name: %w", err)
	}

	return nil
}

// DeletePowerProfile remove um perfil de energia
func (w *WindowsPowerManagementService) DeletePowerProfile(ctx context.Context, profileID string) error {
	cmd := exec.CommandContext(ctx, "powercfg", "/delete", profileID)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete power profile: %w", err)
	}

	return nil
}

// GetPowerProfileSettings retorna as configurações detalhadas de um perfil
func (w *WindowsPowerManagementService) GetPowerProfileSettings(ctx context.Context, profileID string) (map[string]interface{}, error) {
	cmd := exec.CommandContext(ctx, "powercfg", "/query", profileID)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query power profile settings: %w", err)
	}

	// Parse básico da saída - em uma implementação completa, 
	// seria necessário um parser mais sofisticado
	settings := make(map[string]interface{})
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") && !strings.Contains(line, "GUID") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				settings[key] = value
			}
		}
	}

	return settings, nil
}

// SetPowerProfileSetting define uma configuração específica de um perfil
func (w *WindowsPowerManagementService) SetPowerProfileSetting(ctx context.Context, profileID, subgroup, setting, value string) error {
	cmd := exec.CommandContext(ctx, "powercfg", "/setacvalueindex", profileID, subgroup, setting, value)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set power profile setting: %w", err)
	}

	// Aplica as configurações
	cmd = exec.CommandContext(ctx, "powercfg", "/setactive", profileID)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Run()
}
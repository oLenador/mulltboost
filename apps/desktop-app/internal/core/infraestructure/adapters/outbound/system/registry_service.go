//go:build windows
// +build windows

package system

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"golang.org/x/sys/windows/registry"
)

type RegistryKeyInfo struct {
	Path          string
	SubKeyCount   uint32
	ValueCount    uint32
	SubKeys       []string
	ValueNames    []string
	LastWriteTime time.Time
}

type registryService struct{}

// NewRegistryService cria uma nova instância do serviço de registry
func NewRegistryService() outbound.RegistryService {
	return &registryService{}
}

// parseKeyPath separa o root key do sub path
func (rs *registryService) parseKeyPath(keyPath string) (registry.Key, string, error) {
	parts := strings.SplitN(keyPath, "\\", 2)
	if len(parts) < 1 {
		return 0, "", fmt.Errorf("invalid key path: %s", keyPath)
	}

	var rootKey registry.Key
	switch strings.ToUpper(parts[0]) {
	case "HKEY_CLASSES_ROOT", "HKCR":
		rootKey = registry.CLASSES_ROOT
	case "HKEY_CURRENT_USER", "HKCU":
		rootKey = registry.CURRENT_USER
	case "HKEY_LOCAL_MACHINE", "HKLM":
		rootKey = registry.LOCAL_MACHINE
	case "HKEY_USERS", "HKU":
		rootKey = registry.USERS
	case "HKEY_CURRENT_CONFIG", "HKCC":
		rootKey = registry.CURRENT_CONFIG
	default:
		return 0, "", fmt.Errorf("unsupported root key: %s", parts[0])
	}

	subPath := ""
	if len(parts) > 1 {
		subPath = parts[1]
	}

	return rootKey, subPath, nil
}

// ReadRegistryValue lê um valor do registry
func (rs *registryService) ReadRegistryValue(ctx context.Context, keyPath, valueName string) (interface{}, error) {
	rootKey, subPath, err := rs.parseKeyPath(keyPath)
	if err != nil {
		return nil, fmt.Errorf("parsing key path: %w", err)
	}

	key, err := registry.OpenKey(rootKey, subPath, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("opening registry key %s: %w", keyPath, err)
	}
	defer key.Close()

	// Primeiro, vamos descobrir o tipo do valor
	_, valType, err := key.GetValue(valueName, nil)
	if err != nil {
		return nil, fmt.Errorf("getting value type for %s\\%s: %w", keyPath, valueName, err)
	}

	switch valType {
	case registry.SZ, registry.EXPAND_SZ:
		val, _, err := key.GetStringValue(valueName)
		if err != nil {
			return nil, fmt.Errorf("reading string value %s\\%s: %w", keyPath, valueName, err)
		}
		return val, nil

	case registry.DWORD:
		val, _, err := key.GetIntegerValue(valueName)
		if err != nil {
			return nil, fmt.Errorf("reading dword value %s\\%s: %w", keyPath, valueName, err)
		}
		return uint32(val), nil

	case registry.QWORD:
		val, _, err := key.GetIntegerValue(valueName)
		if err != nil {
			return nil, fmt.Errorf("reading qword value %s\\%s: %w", keyPath, valueName, err)
		}
		return val, nil

	case registry.MULTI_SZ:
		val, _, err := key.GetStringsValue(valueName)
		if err != nil {
			return nil, fmt.Errorf("reading multi-string value %s\\%s: %w", keyPath, valueName, err)
		}
		return val, nil

	case registry.BINARY:
		val, _, err := key.GetBinaryValue(valueName)
		if err != nil {
			return nil, fmt.Errorf("reading binary value %s\\%s: %w", keyPath, valueName, err)
		}
		return val, nil

	default:
		return nil, fmt.Errorf("unsupported registry value type: %d", valType)
	}
}

// WriteRegistryValue escreve um valor no registry
func (rs *registryService) WriteRegistryValue(ctx context.Context, keyPath, valueName string, value interface{}) error {
	rootKey, subPath, err := rs.parseKeyPath(keyPath)
	if err != nil {
		return fmt.Errorf("parsing key path: %w", err)
	}

	key, err := registry.OpenKey(rootKey, subPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("opening registry key %s: %w", keyPath, err)
	}
	defer key.Close()

	switch v := value.(type) {
	case string:
		err = key.SetStringValue(valueName, v)
	case uint32:
		err = key.SetDWordValue(valueName, v)
	case uint64:
		err = key.SetQWordValue(valueName, v)
	case []string:
		err = key.SetStringsValue(valueName, v)
	case []byte:
		err = key.SetBinaryValue(valueName, v)
	case int:
		err = key.SetDWordValue(valueName, uint32(v))
	case int32:
		err = key.SetDWordValue(valueName, uint32(v))
	case int64:
		err = key.SetQWordValue(valueName, uint64(v))
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}

	if err != nil {
		return fmt.Errorf("setting registry value %s\\%s: %w", keyPath, valueName, err)
	}

	return nil
}

// DeleteRegistryValue remove um valor do registry
func (rs *registryService) DeleteRegistryValue(ctx context.Context, keyPath, valueName string) error {
	rootKey, subPath, err := rs.parseKeyPath(keyPath)
	if err != nil {
		return fmt.Errorf("parsing key path: %w", err)
	}

	key, err := registry.OpenKey(rootKey, subPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("opening registry key %s: %w", keyPath, err)
	}
	defer key.Close()

	err = key.DeleteValue(valueName)
	if err != nil {
		return fmt.Errorf("deleting registry value %s\\%s: %w", keyPath, valueName, err)
	}

	return nil
}

// CreateRegistryKey cria uma nova chave no registry
func (rs *registryService) CreateRegistryKey(ctx context.Context, keyPath string) error {
	rootKey, subPath, err := rs.parseKeyPath(keyPath)
	if err != nil {
		return fmt.Errorf("parsing key path: %w", err)
	}

	key, _, err := registry.CreateKey(rootKey, subPath, registry.CREATE_SUB_KEY)
	if err != nil {
		return fmt.Errorf("creating registry key %s: %w", keyPath, err)
	}
	defer key.Close()

	return nil
}

// DeleteRegistryKey remove uma chave do registry
func (rs *registryService) DeleteRegistryKey(ctx context.Context, keyPath string) error {
	rootKey, subPath, err := rs.parseKeyPath(keyPath)
	if err != nil {
		return fmt.Errorf("parsing key path: %w", err)
	}

	err = registry.DeleteKey(rootKey, subPath)
	if err != nil {
		return fmt.Errorf("deleting registry key %s: %w", keyPath, err)
	}

	return nil
}

// KeyExists verifica se uma chave existe no registry
func (rs *registryService) KeyExists(ctx context.Context, keyPath string) (bool, error) {
	rootKey, subPath, err := rs.parseKeyPath(keyPath)
	if err != nil {
		return false, fmt.Errorf("parsing key path: %w", err)
	}

	key, err := registry.OpenKey(rootKey, subPath, registry.QUERY_VALUE)
	if err != nil {
		// Se o erro é "key not found", retorna false sem erro
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("checking if key exists %s: %w", keyPath, err)
	}
	defer key.Close()

	return true, nil
}

// BackupRegistryKey faz backup de uma chave do registry
func (rs *registryService) BackupRegistryKey(ctx context.Context, keyPath, backupPath string) error {
	// Usando reg.exe para backup completo (mais confiável que implementação manual)
	cmd := exec.CommandContext(ctx, "reg", "export", keyPath, backupPath, "/y")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("backing up registry key %s to %s: %w (output: %s)", 
			keyPath, backupPath, err, string(output))
	}

	return nil
}

// RestoreRegistryKey restaura uma chave do registry de um backup
func (rs *registryService) RestoreRegistryKey(ctx context.Context, keyPath, backupPath string) error {
	// Usando reg.exe para restore
	cmd := exec.CommandContext(ctx, "reg", "import", backupPath)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("restoring registry key from %s: %w (output: %s)", 
			backupPath, err, string(output))
	}

	return nil
}

// ApplyPerformanceRegistryTweaks aplica otimizações de performance no registry
func (rs *registryService) ApplyPerformanceRegistryTweaks(ctx context.Context) error {
	tweaks := map[string]map[string]interface{}{
		// Otimizações de performance do sistema
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\PriorityControl`: {
			"Win32PrioritySeparation": uint32(2), // Otimizar para aplicações
		},
		
		// Otimizações de rede
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`: {
			"TcpAckFrequency":     uint32(1),
			"TCPNoDelay":          uint32(1),
			"DefaultTTL":          uint32(64),
			"EnablePMTUDiscovery": uint32(1),
		},
		
		// Otimizações de disco
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\FileSystem`: {
			"NtfsDisableLastAccessUpdate": uint32(1), // Desabilita atualização de último acesso
		},
		
		// Otimizações visuais para performance
		`HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`: {
			"VisualFXSetting": uint32(2), // Ajustar para melhor performance
		},
	}

	for keyPath, values := range tweaks {
		// Primeiro, criar a chave se não existir
		if err := rs.CreateRegistryKey(ctx, keyPath); err != nil {
			// Se falhar ao criar, tenta continuar (pode já existir)
		}

		// Aplicar cada valor
		for valueName, value := range values {
			if err := rs.WriteRegistryValue(ctx, keyPath, valueName, value); err != nil {
				return fmt.Errorf("applying performance tweak %s\\%s: %w", keyPath, valueName, err)
			}
		}
	}

	return nil
}

// RestoreDefaultRegistrySettings restaura configurações padrão
func (rs *registryService) RestoreDefaultRegistrySettings(ctx context.Context) error {
	defaults := map[string]map[string]interface{}{
		// Restaurar configurações padrão de prioridade
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\PriorityControl`: {
			"Win32PrioritySeparation": uint32(26), // Valor padrão
		},
		
		// Restaurar configurações padrão de rede
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`: {
			"DefaultTTL": uint32(128), // Valor padrão do Windows
		},
		
		// Restaurar configurações de arquivo
		`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\FileSystem`: {
			"NtfsDisableLastAccessUpdate": uint32(0), // Habilitar novamente
		},
		
		// Restaurar efeitos visuais
		`HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`: {
			"VisualFXSetting": uint32(0), // Deixar Windows decidir
		},
	}

	for keyPath, values := range defaults {
		for valueName, value := range values {
			if err := rs.WriteRegistryValue(ctx, keyPath, valueName, value); err != nil {
				return fmt.Errorf("restoring default setting %s\\%s: %w", keyPath, valueName, err)
			}
		}
	}

	return nil
}
//go:build windows
// +build windows
package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"golang.org/x/sys/windows/registry"
)

type windowsCacheManagementService struct {
	scheduleRegistry string
	tempDirs         []string
	browserPaths     map[string][]string
}

// NewCacheManagementService cria uma nova instância do serviço
func NewCacheManagementService() CacheManagementService {
	return &windowsCacheManagementService{
		scheduleRegistry: `SOFTWARE\MultBoost\CacheManagement`,
		tempDirs: []string{
			os.TempDir(),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Temp"),
			filepath.Join(os.Getenv("WINDIR"), "Temp"),
			filepath.Join(os.Getenv("WINDIR"), "Prefetch"),
		},
		browserPaths: map[string][]string{
			"chrome": {
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data", "Default", "Cache"),
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data", "Default", "Code Cache"),
			},
			"firefox": {
				filepath.Join(os.Getenv("APPDATA"), "Mozilla", "Firefox", "Profiles"),
			},
			"edge": {
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Edge", "User Data", "Default", "Cache"),
			},
		},
	}
}

// ClearSystemCache limpa o cache do sistema usando comandos nativos
func (s *windowsCacheManagementService) ClearSystemCache(ctx context.Context) error {
	commands := []struct {
		name string
		args []string
	}{
		{"sfc", []string{"/scannow"}},
		{"dism", []string{"/online", "/cleanup-image", "/restorehealth"}},
	}

	for _, cmd := range commands {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := s.executeCommand(ctx, cmd.name, cmd.args...); err != nil {
				return fmt.Errorf("failed to execute %s: %w", cmd.name, err)
			}
		}
	}

	return nil
}

// ClearDNSCache limpa o cache DNS do sistema
func (s *windowsCacheManagementService) ClearDNSCache(ctx context.Context) error {
	commands := [][]string{
		{"ipconfig", "/flushdns"},
		{"netsh", "int", "ip", "reset"},
		{"netsh", "winsock", "reset"},
	}

	for _, args := range commands {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := s.executeCommand(ctx, args[0], args[1:]...); err != nil {
				return fmt.Errorf("failed to execute %s: %w", strings.Join(args, " "), err)
			}
		}
	}

	return nil
}

// ClearTemporaryFiles remove arquivos temporários do sistema
func (s *windowsCacheManagementService) ClearTemporaryFiles(ctx context.Context) error {
	var totalErrors []string

	for _, tempDir := range s.tempDirs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := s.clearDirectory(ctx, tempDir); err != nil {
				totalErrors = append(totalErrors, fmt.Sprintf("error clearing %s: %v", tempDir, err))
			}
		}
	}

	// Executa limpeza de disco do Windows
	if err := s.executeCommand(ctx, "cleanmgr", "/sagerun:1"); err != nil {
		totalErrors = append(totalErrors, fmt.Sprintf("cleanmgr error: %v", err))
	}

	if len(totalErrors) > 0 {
		return fmt.Errorf("temporary files cleanup completed with errors: %s", strings.Join(totalErrors, "; "))
	}

	return nil
}

// ClearBrowserCache limpa cache dos navegadores principais
func (s *windowsCacheManagementService) ClearBrowserCache(ctx context.Context) error {
	var totalErrors []string

	for browser, paths := range s.browserPaths {
		for _, path := range paths {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if browser == "firefox" {
					// Firefox precisa de tratamento especial para profiles
					if err := s.clearFirefoxCache(ctx, path); err != nil {
						totalErrors = append(totalErrors, fmt.Sprintf("firefox cache error: %v", err))
					}
				} else {
					if err := s.clearDirectory(ctx, path); err != nil {
						totalErrors = append(totalErrors, fmt.Sprintf("%s cache error: %v", browser, err))
					}
				}
			}
		}
	}

	if len(totalErrors) > 0 {
		return fmt.Errorf("browser cache cleanup completed with errors: %s", strings.Join(totalErrors, "; "))
	}

	return nil
}

// AnalyzeCacheUsage analisa o uso de cache no sistema
func (s *windowsCacheManagementService) AnalyzeCacheUsage(ctx context.Context) (int64, error) {
	var totalSize int64

	// Analisa diretórios temporários
	for _, tempDir := range s.tempDirs {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			size, err := s.getDirSize(tempDir)
			if err == nil {
				totalSize += size
			}
		}
	}

	// Analisa cache dos navegadores
	for _, paths := range s.browserPaths {
		for _, path := range paths {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			default:
				size, err := s.getDirSize(path)
				if err == nil {
					totalSize += size
				}
			}
		}
	}

	return totalSize, nil
}

// GetCacheStatistics retorna estatísticas detalhadas do cache
func (s *windowsCacheManagementService) GetCacheStatistics(ctx context.Context) (*entities.CacheData, error) {
	stats := &entities.CacheData{
		TotalSize:    0,
		TempFiles:    make(map[string]int64),
		BrowserCache: make(map[string]int64),
		SystemCache:  0,
		LastCleanup:  time.Now(),
	}

	// Estatísticas de arquivos temporários
	for _, tempDir := range s.tempDirs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			size, err := s.getDirSize(tempDir)
			if err == nil {
				stats.TempFiles[filepath.Base(tempDir)] = size
				stats.TotalSize += size
			}
		}
	}

	// Estatísticas de cache de navegadores
	for browser, paths := range s.browserPaths {
		var browserTotal int64
		for _, path := range paths {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				size, err := s.getDirSize(path)
				if err == nil {
					browserTotal += size
				}
			}
		}
		stats.BrowserCache[browser] = browserTotal
		stats.TotalSize += browserTotal
	}

	// Sistema de cache (estimativa baseada em RAM e paging file)
	systemSize, err := s.getSystemCacheSize()
	if err == nil {
		stats.SystemCache = systemSize
		stats.TotalSize += systemSize
	}

	return stats, nil
}

// SetCacheCleanupSchedule configura agendamento de limpeza
func (s *windowsCacheManagementService) SetCacheCleanupSchedule(ctx context.Context, schedule string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, s.scheduleRegistry, registry.CREATE_SUB_KEY|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue("CleanupSchedule", schedule); err != nil {
		return fmt.Errorf("failed to set schedule: %w", err)
	}

	// Cria tarefa no Task Scheduler do Windows
	taskName := "MultBoostCacheCleanup"
	taskCmd := fmt.Sprintf(`schtasks /create /tn "%s" /tr "multboost.exe --cleanup-cache" /sc %s /f`, 
		taskName, schedule)
	
	if err := s.executeCommand(ctx, "cmd", "/c", taskCmd); err != nil {
		return fmt.Errorf("failed to create scheduled task: %w", err)
	}

	return nil
}

// DisableAutomaticCacheCleanup desabilita limpeza automática
func (s *windowsCacheManagementService) DisableAutomaticCacheCleanup(ctx context.Context) error {
	// Remove do registry
	key, err := registry.OpenKey(registry.CURRENT_USER, s.scheduleRegistry, registry.SET_VALUE)
	if err == nil {
		key.DeleteValue("CleanupSchedule")
		key.Close()
	}

	// Remove tarefa do Task Scheduler
	taskName := "MultBoostCacheCleanup"
	taskCmd := fmt.Sprintf(`schtasks /delete /tn "%s" /f`, taskName)
	
	if err := s.executeCommand(ctx, "cmd", "/c", taskCmd); err != nil {
		return fmt.Errorf("failed to delete scheduled task: %w", err)
	}

	return nil
}

// Métodos auxiliares privados

func (s *windowsCacheManagementService) executeCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}
	return nil
}

func (s *windowsCacheManagementService) clearDirectory(ctx context.Context, dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil // Diretório não existe, não é erro
	}

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continua mesmo com erros individuais
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if !info.IsDir() && path != dirPath {
				os.Remove(path) // Ignora erros de arquivos em uso
			}
		}
		return nil
	})
}

func (s *windowsCacheManagementService) clearFirefoxCache(ctx context.Context, profilesDir string) error {
	if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(profilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if info.IsDir() && strings.Contains(info.Name(), ".default") {
				cacheDir := filepath.Join(path, "cache2")
				if err := s.clearDirectory(ctx, cacheDir); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *windowsCacheManagementService) getDirSize(dirPath string) (int64, error) {
	var size int64
	
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continua mesmo com erros
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	
	return size, err
}

func (s *windowsCacheManagementService) getSystemCacheSize() (int64, error) {
	// Usando syscalls do Windows para obter informações de memória
	var memStatus struct {
		Length               uint32
		MemoryLoad           uint32
		TotalPhys           uint64
		AvailPhys           uint64
		TotalPageFile       uint64
		AvailPageFile       uint64
		TotalVirtual        uint64
		AvailVirtual        uint64
		AvailExtendedVirtual uint64
	}

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx := kernel32.NewProc("GlobalMemoryStatusEx")

	memStatus.Length = uint32(unsafe.Sizeof(memStatus))
	ret, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	
	if ret == 0 {
		return 0, fmt.Errorf("failed to get memory status")
	}

	// Estimativa do cache do sistema (diferença entre total e disponível)
	systemCache := int64(memStatus.TotalPhys - memStatus.AvailPhys)
	return systemCache, nil
}
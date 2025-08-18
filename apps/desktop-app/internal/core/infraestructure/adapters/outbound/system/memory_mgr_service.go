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

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"golang.org/x/sys/windows"
)

// Windows API constants and structures
const (
	SYSTEM_CACHE_INFORMATION = 21
	SYSTEM_FILECACHE_INFORMATION = 45
)

type SYSTEM_CACHE_INFORMATION_STRUCT struct {
	CurrentSize                     uint64
	PeakSize                        uint64
	PageFaultCount                  uint32
	MinimumWorkingSet              uint64
	MaximumWorkingSet              uint64
	CurrentSizeIncludingTransitionInPages uint64
	PeakSizeIncludingTransitionInPages    uint64
	TransitionRePurposeCount             uint32
	Flags                               uint32
}

// WindowsMemoryManagementService implements MemoryManagementService
type WindowsMemoryManagementService struct {
	wmiExecutor       *WMIExecutor
	powershellExecutor *PowerShellExecutor
	systemCallExecutor *SystemCallExecutor
}

// WMIExecutor handles WMI operations
type WMIExecutor struct{}

// PowerShellExecutor handles PowerShell operations  
type PowerShellExecutor struct{}

// SystemCallExecutor handles direct Windows API calls
type SystemCallExecutor struct{}

// NewWindowsMemoryManagementService creates a new instance
func NewWindowsMemoryManagementService() *WindowsMemoryManagementService {
	return &WindowsMemoryManagementService{
		wmiExecutor:       &WMIExecutor{},
		powershellExecutor: &PowerShellExecutor{},
		systemCallExecutor: &SystemCallExecutor{},
	}
}

// ClearMemoryCache clears system memory cache
func (s *WindowsMemoryManagementService) ClearMemoryCache(ctx context.Context) error {
	// Use direct system calls for performance-critical operations
	if err := s.systemCallExecutor.ClearSystemCache(ctx); err != nil {
		// Fallback to PowerShell method
		return s.powershellExecutor.ClearMemoryCache(ctx)
	}
	return nil
}

// OptimizeMemoryUsage optimizes memory usage across the system
func (s *WindowsMemoryManagementService) OptimizeMemoryUsage(ctx context.Context) error {
	// Multi-step optimization process
	steps := []func(context.Context) error{
		s.systemCallExecutor.TrimWorkingSets,
		s.systemCallExecutor.EmptyWorkingSets,
		s.powershellExecutor.RunMemoryDiagnostic,
		s.systemCallExecutor.FlushMemoryCache,
	}

	for i, step := range steps {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := step(ctx); err != nil {
				return fmt.Errorf("memory optimization failed at step %d: %w", i+1, err)
			}
		}
	}

	return nil
}

// SetVirtualMemorySize configures virtual memory (page file) size
func (s *WindowsMemoryManagementService) SetVirtualMemorySize(ctx context.Context, sizeGB int) error {
	return s.powershellExecutor.SetVirtualMemorySize(ctx, sizeGB)
}

// GetMemoryInfo retrieves comprehensive memory information
func (s *WindowsMemoryManagementService) GetMemoryInfo(ctx context.Context) (*entities.MemoryInfo, error) {
	return s.wmiExecutor.GetMemoryInfo(ctx)
}

// GetMemoryUsageByProcess gets memory usage for each process
func (s *WindowsMemoryManagementService) GetMemoryUsageByProcess(ctx context.Context) (map[string]float64, error) {
	return s.wmiExecutor.GetMemoryUsageByProcess(ctx)
}

// EnableMemoryCompression enables memory compression feature
func (s *WindowsMemoryManagementService) EnableMemoryCompression(ctx context.Context) error {
	return s.powershellExecutor.SetMemoryCompression(ctx, true)
}

// DisableMemoryCompression disables memory compression feature
func (s *WindowsMemoryManagementService) DisableMemoryCompression(ctx context.Context) error {
	return s.powershellExecutor.SetMemoryCompression(ctx, false)
}

// SetPageFileLocation sets the page file location
func (s *WindowsMemoryManagementService) SetPageFileLocation(ctx context.Context, drive string) error {
	return s.powershellExecutor.SetPageFileLocation(ctx, drive)
}

// WMIExecutor implementations
func (w *WMIExecutor) GetMemoryInfo(ctx context.Context) (*entities.MemoryInfo, error) {
	query := `wmic computersystem get TotalPhysicalMemory /format:list`
	
	cmd := exec.CommandContext(ctx, "cmd", "/C", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get total physical memory: %w", err)
	}

	totalMemory, err := parseWMIOutput(string(output), "TotalPhysicalMemory")
	if err != nil {
		return nil, fmt.Errorf("failed to parse total memory: %w", err)
	}

	// Get available memory
	availQuery := `wmic OS get AvailablePhysicalMemory /format:list`
	cmd = exec.CommandContext(ctx, "cmd", "/C", availQuery)
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get available physical memory: %w", err)
	}

	availableMemory, err := parseWMIOutput(string(output), "AvailablePhysicalMemory")
	if err != nil {
		return nil, fmt.Errorf("failed to parse available memory: %w", err)
	}

	// Get page file usage
	pageFileQuery := `wmic pagefileset get AllocatedBaseSize,CurrentUsage /format:list`
	cmd = exec.CommandContext(ctx, "cmd", "/C", pageFileQuery)
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get page file info: %w", err)
	}

	pageFileSize, _ := parseWMIOutput(string(output), "AllocatedBaseSize")
	pageFileUsed, _ := parseWMIOutput(string(output), "CurrentUsage")

	usedMemory := totalMemory - availableMemory
	memoryUsagePercent := float64(usedMemory) / float64(totalMemory) * 100

	return &entities.MemoryInfo{
		TotalMemoryGB:      float64(totalMemory) / (1024 * 1024 * 1024),
		AvailableMemoryGB:  float64(availableMemory) / (1024 * 1024 * 1024),
		UsedMemoryGB:       float64(usedMemory) / (1024 * 1024 * 1024),
		MemoryUsagePercent: memoryUsagePercent,
		PageFileSize:       float64(pageFileSize),
		PageFileUsed:       float64(pageFileUsed),
		CacheSize:          0, // Will be populated by system call
	}, nil
}

func (w *WMIExecutor) GetMemoryUsageByProcess(ctx context.Context) (map[string]float64, error) {
	query := `wmic process get Name,WorkingSetSize /format:csv`
	
	cmd := exec.CommandContext(ctx, "cmd", "/C", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get process memory usage: %w", err)
	}

	return parseProcessMemoryOutput(string(output))
}

// PowerShellExecutor implementations
func (p *PowerShellExecutor) ClearMemoryCache(ctx context.Context) error {
	script := `
	[System.GC]::Collect()
	[System.GC]::WaitForPendingFinalizers()
	[System.GC]::Collect()
	`
	
	return p.executeScript(ctx, script)
}

func (p *PowerShellExecutor) SetVirtualMemorySize(ctx context.Context, sizeGB int) error {
	sizeMB := sizeGB * 1024
	script := fmt.Sprintf(`
	$cs = Get-WmiObject -Class Win32_ComputerSystem -EnableAllPrivileges
	$cs.AutomaticManagedPagefile = $False
	$cs.Put()
	
	$pf = Get-WmiObject -Class Win32_PageFileSetting
	if ($pf) {
		$pf.Delete()
	}
	
	Set-WmiInstance -Class Win32_PageFileSetting -Arguments @{
		name = "C:\pagefile.sys"
		InitialSize = %d
		MaximumSize = %d
	}
	`, sizeMB, sizeMB*2)

	return p.executeScript(ctx, script)
}

func (p *PowerShellExecutor) SetMemoryCompression(ctx context.Context, enabled bool) error {
	action := "Disable"
	if enabled {
		action = "Enable"
	}

	script := fmt.Sprintf(`%s-MMAgent -MemoryCompression`, action)
	return p.executeScript(ctx, script)
}

func (p *PowerShellExecutor) SetPageFileLocation(ctx context.Context, drive string) error {
	script := fmt.Sprintf(`
	$pf = Get-WmiObject -Class Win32_PageFileSetting
	if ($pf) {
		$pf.Delete()
	}
	
	Set-WmiInstance -Class Win32_PageFileSetting -Arguments @{
		name = "%s:\pagefile.sys"
		InitialSize = 0
		MaximumSize = 0
	}
	`, drive)

	return p.executeScript(ctx, script)
}

func (p *PowerShellExecutor) RunMemoryDiagnostic(ctx context.Context) error {
	script := `
	Get-Process | Where-Object {$_.WorkingSet -gt 100MB} | 
	ForEach-Object {
		try {
			$_.ProcessorAffinity = $_.ProcessorAffinity
		} catch {}
	}
	`
	return p.executeScript(ctx, script)
}

func (p *PowerShellExecutor) executeScript(ctx context.Context, script string) error {
	cmd := exec.CommandContext(ctx, "powershell", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("powershell execution failed: %w, output: %s", err, string(output))
	}
	
	return nil
}

// SystemCallExecutor implementations
func (s *SystemCallExecutor) ClearSystemCache(ctx context.Context) error {
	// Load ntdll.dll
	ntdll := syscall.NewLazyDLL("ntdll.dll")
	ntSetSystemInformation := ntdll.NewProc("NtSetSystemInformation")

	// Call NtSetSystemInformation to clear cache
	ret, _, err := ntSetSystemInformation.Call(
		uintptr(SYSTEM_CACHE_INFORMATION),
		0,
		0,
	)

	if ret != 0 {
		return fmt.Errorf("failed to clear system cache: %w", err)
	}

	return nil
}

func (s *SystemCallExecutor) TrimWorkingSets(ctx context.Context) error {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setProcessWorkingSetSize := kernel32.NewProc("SetProcessWorkingSetSize")

	currentProcess := windows.CurrentProcess()
	
	ret, _, err := setProcessWorkingSetSize.Call(
		uintptr(currentProcess),
		uintptr(0xFFFFFFFF), // Trim working set
		uintptr(0xFFFFFFFF),
	)

	if ret == 0 {
		return fmt.Errorf("failed to trim working sets: %w", err)
	}

	return nil
}

func (s *SystemCallExecutor) EmptyWorkingSets(ctx context.Context) error {
	psapi := syscall.NewLazyDLL("psapi.dll")
	emptyWorkingSet := psapi.NewProc("EmptyWorkingSet")

	currentProcess := windows.CurrentProcess()
	
	ret, _, err := emptyWorkingSet.Call(uintptr(currentProcess))
	if ret == 0 {
		return fmt.Errorf("failed to empty working set: %w", err)
	}

	return nil
}

func (s *SystemCallExecutor) FlushMemoryCache(ctx context.Context) error {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setSystemFileCacheSize := kernel32.NewProc("SetSystemFileCacheSize")

	ret, _, err := setSystemFileCacheSize.Call(
		uintptr(0xFFFFFFFF), // Clear file cache
		uintptr(0xFFFFFFFF),
		0,
	)

	if ret == 0 {
		return fmt.Errorf("failed to flush memory cache: %w", err)
	}

	return nil
}

// Utility functions
func parseWMIOutput(output, field string) (uint64, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, field) {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				if value != "" {
					return strconv.ParseUint(value, 10, 64)
				}
			}
		}
	}
	return 0, fmt.Errorf("field %s not found in WMI output", field)
}

func parseProcessMemoryOutput(output string) (map[string]float64, error) {
	result := make(map[string]float64)
	lines := strings.Split(output, "\n")
	
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // Skip header and empty lines
		}
		
		parts := strings.Split(line, ",")
		if len(parts) >= 3 {
			processName := strings.TrimSpace(parts[1])
			workingSetStr := strings.TrimSpace(parts[2])
			
			if workingSet, err := strconv.ParseUint(workingSetStr, 10, 64); err == nil {
				// Convert to MB
				result[processName] = float64(workingSet) / (1024 * 1024)
			}
		}
	}
	
	return result, nil
}

// Error types for better error handling
type MemoryManagementError struct {
	Operation string
	Err       error
}

func (e *MemoryManagementError) Error() string {
	return fmt.Sprintf("memory management operation '%s' failed: %v", e.Operation, e.Err)
}

func (e *MemoryManagementError) Unwrap() error {
	return e.Err
}

// Wrap errors with operation context
func wrapError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return &MemoryManagementError{
		Operation: operation,
		Err:       err,
	}
}
// internal/infrastructure/services/windows/config.go
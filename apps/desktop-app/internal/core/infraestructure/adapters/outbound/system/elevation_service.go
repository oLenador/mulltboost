//go:build windows
// +build windows
package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowsElevationService struct {
	processHandle windows.Handle
	tokenHandle   windows.Handle
}

// Creates a new WindowsElevationService instance
func NewWindowsElevationService() (*WindowsElevationService, error) {
	service := &WindowsElevationService{}
	var err error
	service.processHandle, err = windows.GetCurrentProcess()
	if err != nil {
		return nil, fmt.Errorf("failed to get current process handle: %w", err)
	}
	return service, nil
}

// Checks if the current process is running with admin privileges
func (w *WindowsElevationService) IsElevated(ctx context.Context) (bool, error) {
	var token windows.Token
	err := windows.OpenProcessToken(w.processHandle, windows.TOKEN_QUERY, &token)
	if err != nil {
		return false, fmt.Errorf("failed to open process token: %w", err)
	}
	defer token.Close()

	var elevation windows.TokenElevation
	var returnedLen uint32
	err = windows.GetTokenInformation(
		token,
		windows.TokenElevation,
		(*byte)(unsafe.Pointer(&elevation)),
		uint32(unsafe.Sizeof(elevation)),
		&returnedLen,
	)
	if err != nil {
		return false, fmt.Errorf("failed to get token elevation info: %w", err)
	}

	return elevation.TokenIsElevated != 0, nil
}

// Checks if a specific operation requires elevation
func (w *WindowsElevationService) RequiresElevation(ctx context.Context, operation string) (bool, error) {
	elevatedOperations := map[string]bool{
		"service_install":   true,
		"service_uninstall": true,
		"service_start":     true,
		"service_stop":      true,
		"registry_write":    true,
		"system_config":     true,
		"driver_install":    true,
		"firewall_config":   true,
		"user_management":   true,
		"file_system_root":  true,
		"network_config":    true,
	}

	requiresElevation, exists := elevatedOperations[operation]
	if !exists {
		return false, nil
	}

	if requiresElevation {
		isElevated, err := w.IsElevated(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to check elevation status: %w", err)
		}
		return !isElevated, nil
	}

	return false, nil
}

// Requests elevation for the current process
func (w *WindowsElevationService) RequestElevation(ctx context.Context) error {
	isElevated, err := w.IsElevated(ctx)
	if err != nil {
		return fmt.Errorf("failed to check current elevation status: %w", err)
	}
	if isElevated {
		return nil
	}

	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	executableUTF16, err := windows.UTF16PtrFromString(executable)
	if err != nil {
		return fmt.Errorf("failed to convert executable path to UTF16: %w", err)
	}

	verbUTF16, err := windows.UTF16PtrFromString("runas")
	if err != nil {
		return fmt.Errorf("failed to convert verb to UTF16: %w", err)
	}

	sei := windows.ShellExecuteInfo{
		CbSize:   uint32(unsafe.Sizeof(windows.ShellExecuteInfo{})),
		FMask:    windows.SEE_MASK_NOCLOSEPROCESS | windows.SEE_MASK_NOASYNC,
		Verb:     verbUTF16,
		File:     executableUTF16,
		Show:     windows.SW_SHOWNORMAL,
		HInstApp: 0,
	}

	err = windows.ShellExecuteEx(&sei)
	if err != nil {
		return fmt.Errorf("failed to request elevation: %w", err)
	}

	os.Exit(0)
	return nil
}

// Runs a command with admin privileges
func (w *WindowsElevationService) RunAsAdmin(ctx context.Context, command string, args []string) error {
	uacEnabled, err := w.IsUACEnabled(ctx)
	if err != nil {
		return fmt.Errorf("failed to check UAC status: %w", err)
	}

	if !uacEnabled {
		cmd := exec.CommandContext(ctx, command, args...)
		return cmd.Run()
	}

	allArgs := strings.Join(args, " ")

	verbUTF16, err := windows.UTF16PtrFromString("runas")
	if err != nil {
		return fmt.Errorf("failed to convert verb to UTF16: %w", err)
	}

	fileUTF16, err := windows.UTF16PtrFromString(command)
	if err != nil {
		return fmt.Errorf("failed to convert command to UTF16: %w", err)
	}

	var parametersUTF16 *uint16
	if allArgs != "" {
		parametersUTF16, err = windows.UTF16PtrFromString(allArgs)
		if err != nil {
			return fmt.Errorf("failed to convert parameters to UTF16: %w", err)
		}
	}

	sei := windows.ShellExecuteInfo{
		CbSize:     uint32(unsafe.Sizeof(windows.ShellExecuteInfo{})),
		FMask:      windows.SEE_MASK_NOCLOSEPROCESS | windows.SEE_MASK_NOASYNC,
		Verb:       verbUTF16,
		File:       fileUTF16,
		Parameters: parametersUTF16,
		Show:       windows.SW_HIDE,
		HInstApp:   0,
	}

	err = windows.ShellExecuteEx(&sei)
	if err != nil {
		return fmt.Errorf("failed to run command as admin: %w", err)
	}

	if sei.HProcess != 0 {
		defer windows.CloseHandle(sei.HProcess)
		event, err := windows.WaitForSingleObject(sei.HProcess, windows.INFINITE)
		if err != nil {
			return fmt.Errorf("failed to wait for process completion: %w", err)
		}

		if event != windows.WAIT_OBJECT_0 {
			return fmt.Errorf("process wait returned unexpected result: %d", event)
		}

		var exitCode uint32
		err = windows.GetExitCodeProcess(sei.HProcess, &exitCode)
		if err != nil {
			return fmt.Errorf("failed to get process exit code: %w", err)
		}

		if exitCode != 0 {
			return fmt.Errorf("command failed with exit code: %d", exitCode)
		}
	}

	return nil
}

// Checks if UAC is enabled
func (w *WindowsElevationService) IsUACEnabled(ctx context.Context) (bool, error) {
	key, err := windows.OpenKey(
		windows.HKEY_LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`,
		windows.KEY_READ,
	)
	if err != nil {
		return false, fmt.Errorf("failed to open UAC registry key: %w", err)
	}
	defer windows.CloseKey(key)

	var enableLUA uint32
	var bufSize uint32 = 4
	err = windows.GetValue(key, "EnableLUA", (*byte)(unsafe.Pointer(&enableLUA)), &bufSize)
	if err != nil {
		if err == windows.ERROR_FILE_NOT_FOUND {
			return true, nil
		}
		return false, fmt.Errorf("failed to read EnableLUA value: %w", err)
	}

	return enableLUA != 0, nil
}

// Enables or disables UAC
func (w *WindowsElevationService) ConfigureUAC(ctx context.Context, enabled bool) error {
	isElevated, err := w.IsElevated(ctx)
	if err != nil {
		return fmt.Errorf("failed to check elevation status: %w", err)
	}

	if !isElevated {
		return fmt.Errorf("configuring UAC requires administrative privileges")
	}

	key, err := windows.OpenKey(
		windows.HKEY_LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`,
		windows.KEY_WRITE,
	)
	if err != nil {
		return fmt.Errorf("failed to open UAC registry key for writing: %w", err)
	}
	defer windows.CloseKey(key)

	var value uint32 = 0
	if enabled {
		value = 1
	}

	err = windows.SetValue(
		key,
		"EnableLUA",
		windows.REG_DWORD,
		(*byte)(unsafe.Pointer(&value)),
		4,
	)
	if err != nil {
		return fmt.Errorf("failed to set EnableLUA value: %w", err)
	}

	return nil
}

// Closes resources used by the service
func (w *WindowsElevationService) Close() error {
	return nil
}

// Alias for IsElevated
func (w *WindowsElevationService) IsAdmin(ctx context.Context) (bool, error) {
	return w.IsElevated(ctx)
}

// Alias for RunAsAdmin
func (w *WindowsElevationService) RunElevated(ctx context.Context, command string, args []string) error {
	return w.RunAsAdmin(ctx, command, args)
}

// Returns detailed elevation info
func (w *WindowsElevationService) GetElevationInfo(ctx context.Context) (*ElevationInfo, error) {
	isElevated, err := w.IsElevated(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check elevation status: %w", err)
	}

	uacEnabled, err := w.IsUACEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check UAC status: %w", err)
	}

	return &ElevationInfo{
		IsElevated: isElevated,
		UACEnabled: uacEnabled,
		ProcessID:  uint32(os.Getpid()),
	}, nil
}

type ElevationInfo struct {
	IsElevated bool   `json:"is_elevated"`
	UACEnabled bool   `json:"uac_enabled"`
	ProcessID  uint32 `json:"process_id"`
}

// String representation of ElevationInfo
func (ei *ElevationInfo) String() string {
	return fmt.Sprintf("ElevationInfo{IsElevated: %t, UACEnabled: %t, ProcessID: %d}",
		ei.IsElevated, ei.UACEnabled, ei.ProcessID)
}

const (
	OpServiceInstall   = "service_install"
	OpServiceUninstall = "service_uninstall"
	OpServiceStart     = "service_start"
	OpServiceStop      = "service_stop"
	OpRegistryWrite    = "registry_write"
	OpSystemConfig     = "system_config"
	OpDriverInstall    = "driver_install"
	OpFirewallConfig   = "firewall_config"
	OpUserManagement   = "user_management"
	OpFileSystemRoot   = "file_system_root"
	OpNetworkConfig    = "network_config"
)

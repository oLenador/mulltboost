//go:build windows
// +build windows
package system

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/sys/windows"
)

// WindowsSystemService implementa SystemAPIService para Windows
type WindowsSystemService struct {
	loadedDLLs map[string]*syscall.DLL
	mu         sync.RWMutex
}

// NewWindowsSystemService cria uma nova instância do serviço
func NewWindowsSystemService() *WindowsSystemService {
	return &WindowsSystemService{
		loadedDLLs: make(map[string]*syscall.DLL),
	}
}

// CallWindowsAPI chama uma API específica do Windows
func (w *WindowsSystemService) CallWindowsAPI(ctx context.Context, apiName string, params []interface{}) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	switch apiName {
	case "GetComputerName":
		return w.getComputerName()
	case "GetUserName":
		return w.getUserName()
	case "GetSystemDirectory":
		return w.getSystemDirectory()
	case "GetWindowsDirectory":
		return w.getWindowsDirectory()
	case "GetTempPath":
		return w.getTempPath()
	case "GetLogicalDrives":
		return w.getLogicalDrives()
	default:
		return nil, fmt.Errorf("API não suportada: %s", apiName)
	}
}

// LoadSystemDLL carrega uma DLL do sistema
func (w *WindowsSystemService) LoadSystemDLL(ctx context.Context, dllName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.loadedDLLs[dllName]; exists {
		return nil // Já carregada
	}

	dll, err := syscall.LoadDLL(dllName)
	if err != nil {
		return fmt.Errorf("falha ao carregar DLL %s: %w", dllName, err)
	}

	w.loadedDLLs[dllName] = dll
	return nil
}

// UnloadSystemDLL descarrega uma DLL do sistema
func (w *WindowsSystemService) UnloadSystemDLL(ctx context.Context, dllName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	dll, exists := w.loadedDLLs[dllName]
	if !exists {
		return fmt.Errorf("DLL %s não está carregada", dllName)
	}

	err := dll.Release()
	if err != nil {
		return fmt.Errorf("falha ao descarregar DLL %s: %w", dllName, err)
	}

	delete(w.loadedDLLs, dllName)
	return nil
}

// GetSystemVersion retorna a versão do sistema
func (w *WindowsSystemService) GetSystemVersion(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Usando RtlGetVersion para obter versão real (não afetado por compatibility shims)
	version := windows.RtlGetVersion()

	return fmt.Sprintf("%d.%d.%d", version.MajorVersion, version.MinorVersion, version.BuildNumber), nil
}

// GetSystemArchitecture retorna a arquitetura do sistema
func (w *WindowsSystemService) GetSystemArchitecture(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	return runtime.GOARCH, nil
}

// IsWindows11 verifica se é Windows 11
func (w *WindowsSystemService) IsWindows11(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	version := windows.RtlGetVersion()

	// Windows 11 é identificado pela build number >= 22000
	return version.BuildNumber >= 22000, nil
}

// ExecutePowerShellCommand executa um comando PowerShell
func (w *WindowsSystemService) ExecutePowerShellCommand(ctx context.Context, command string) (string, error) {
	cmd := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: windows.CREATE_NO_WINDOW,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("falha ao executar PowerShell: %w\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

// ExecuteCMDCommand executa um comando CMD
func (w *WindowsSystemService) ExecuteCMDCommand(ctx context.Context, command string) (string, error) {
	cmd := exec.CommandContext(ctx, "cmd.exe", "/C", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: windows.CREATE_NO_WINDOW,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("falha ao executar CMD: %w\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

// RunSystemCommand executa um comando do sistema
func (w *WindowsSystemService) RunSystemCommand(ctx context.Context, command string, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: windows.CREATE_NO_WINDOW,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("falha ao executar comando %s: %w\nOutput: %s", command, err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

// Métodos auxiliares para APIs específicas do Windows

func (w *WindowsSystemService) getComputerName() (string, error) {
	var size uint32 = windows.MAX_COMPUTERNAME_LENGTH + 1
	buf := make([]uint16, size)
	
	err := windows.GetComputerName(&buf[0], &size)
	if err != nil {
		return "", err
	}
	
	return windows.UTF16ToString(buf[:size]), nil
}

func (w *WindowsSystemService) getUserName() (string, error) {
	var size uint32 = 256
	buf := make([]uint16, size)

	err := windows.GetUserNameEx(windows.NameSamCompatible, &buf[0], &size)
	if err != nil {
		return "", err
	}

	return windows.UTF16ToString(buf[:size]), nil
}

func (w *WindowsSystemService) getSystemDirectory() (string, error) {
	dir, err := windows.GetSystemDirectory()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (w *WindowsSystemService) getWindowsDirectory() (string, error) {
	dir, err := windows.GetWindowsDirectory()
	if err != nil {
		return "", err
	}
	return dir, nil
}


func (w *WindowsSystemService) getTempPath() (string, error) {
	var size uint32 = windows.MAX_PATH
	buf := make([]uint16, size)

	n, err := windows.GetTempPath(size, &buf[0])
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", fmt.Errorf("falha ao obter path temporário")
	}

	return windows.UTF16ToString(buf[:n]), nil
}


func (w *WindowsSystemService) getLogicalDrives() ([]string, error) {
	drives, err := windows.GetLogicalDrives()
	if err != nil {
		return nil, err
	}

	var result []string
	for i := 0; i < 26; i++ {
		if drives&(1<<uint(i)) != 0 {
			drive := string(rune('A'+i)) + ":"
			result = append(result, drive)
		}
	}
	return result, nil
}

// Cleanup libera recursos
func (w *WindowsSystemService) Cleanup() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var errors []string
	for name, dll := range w.loadedDLLs {
		if err := dll.Release(); err != nil {
			errors = append(errors, fmt.Sprintf("erro ao liberar %s: %v", name, err))
		}
	}

	w.loadedDLLs = make(map[string]*syscall.DLL)

	if len(errors) > 0 {
		return fmt.Errorf("erros durante cleanup: %s", strings.Join(errors, "; "))
	}

	return nil
}

// SystemInfo retorna informações detalhadas do sistema
func (w *WindowsSystemService) SystemInfo(ctx context.Context) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	version, err := w.GetSystemVersion(ctx)
	if err == nil {
		info["version"] = version
	}

	arch, err := w.GetSystemArchitecture(ctx)
	if err == nil {
		info["architecture"] = arch
	}

	isWin11, err := w.IsWindows11(ctx)
	if err == nil {
		info["isWindows11"] = isWin11
	}

	computerName, err := w.getComputerName()
	if err == nil {
		info["computerName"] = computerName
	}

	userName, err := w.getUserName()
	if err == nil {
		info["userName"] = userName
	}

	drives, err := w.getLogicalDrives()
	if err == nil {
		info["logicalDrives"] = drives
	}

	return info, nil
}

// Constantes para GetSystemMetrics
const (
	SM_CXSCREEN      = 0  // Largura da tela primária
	SM_CYSCREEN      = 1  // Altura da tela primária
	SM_CXVSCROLL     = 2  // Largura da scroll bar vertical
	SM_CYHSCROLL     = 3  // Altura da scroll bar horizontal
	SM_CYCAPTION     = 4  // Altura da barra de título
	SM_CXBORDER      = 5  // Largura da borda
	SM_CYBORDER      = 6  // Altura da borda
	SM_CXFIXEDFRAME  = 7  // Largura do frame fixo
	SM_CYFIXEDFRAME  = 8  // Altura do frame fixo
	SM_CYVTHUMB      = 9  // Altura do thumb vertical
	SM_CXHTHUMB      = 10 // Largura do thumb horizontal
	SM_CXICON        = 11 // Largura padrão do ícone
	SM_CYICON        = 12 // Altura padrão do ícone
	SM_CXCURSOR      = 13 // Largura do cursor
	SM_CYCURSOR      = 14 // Altura do cursor
	SM_CYMENU        = 15 // Altura do menu
	SM_CXFULLSCREEN  = 16 // Largura da área cliente fullscreen
	SM_CYFULLSCREEN  = 17 // Altura da área cliente fullscreen
	SM_CYKANJIWINDOW = 18 // Altura da janela Kanji
	SM_MOUSEPRESENT  = 19 // Mouse presente
	SM_CYVSCROLL     = 20 // Altura da seta da scroll bar vertical
	SM_CXHSCROLL     = 21 // Largura da seta da scroll bar horizontal
	SM_DEBUG         = 22 // Versão debug do Windows
	SM_SWAPBUTTON    = 23 // Botões do mouse trocados
	SM_CXMIN         = 28 // Largura mínima da janela
	SM_CYMIN         = 29 // Altura mínima da janela
	SM_CXSIZE        = 30 // Largura do botão de redimensionar
	SM_CYSIZE        = 31 // Altura do botão de redimensionar
)

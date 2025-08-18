//go:build windows
// +build windows

package system

import (
	"context"
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	outbound "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

// Constantes de prioridade do Windows
const (
	IDLE_PRIORITY_CLASS         = 0x00000040
	BELOW_NORMAL_PRIORITY_CLASS = 0x00004000
	NORMAL_PRIORITY_CLASS       = 0x00000020
	ABOVE_NORMAL_PRIORITY_CLASS = 0x00008000
	HIGH_PRIORITY_CLASS         = 0x00000080
	REALTIME_PRIORITY_CLASS     = 0x00000100
)

// Constantes de acesso a processo
const (
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_SET_INFORMATION   = 0x0200
	PROCESS_TERMINATE         = 0x0001
	PROCESS_VM_READ           = 0x0010
	PROCESS_ALL_ACCESS        = windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0xFFFF
)

// Estrutura para informações do processo (para PROCESSENTRY32)
type processEntry32 struct {
	Size              uint32
	Usage             uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	Threads           uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [260]uint16
}

type windowsProcessPriorityService struct {
	kernel32 *windows.LazyDLL
}

// NewProcessPriorityService cria uma nova instância do serviço
func NewProcessPriorityService() outbound.ProcessPriorityService {
	return &windowsProcessPriorityService{
		kernel32: windows.NewLazySystemDLL("kernel32.dll"),
	}
}

// SetProcessPriority define a prioridade de um processo pelo nome
func (s *windowsProcessPriorityService) SetProcessPriority(ctx context.Context, processName string, priority int) error {
	// Encontrar o processo pelo nome
	processes, err := s.GetRunningProcesses(ctx)
	if err != nil {
		return fmt.Errorf("erro ao buscar processos: %w", err)
	}

	var targetProcessID int
	found := false
	for _, proc := range processes {
		if strings.EqualFold(proc.Name, processName) || strings.EqualFold(proc.Path, processName) {
			targetProcessID = proc.PID
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("processo '%s' não encontrado", processName)
	}

	// Converter prioridade do nosso sistema (0-5) para Windows
	windowsPriority := s.convertToWindowsPriority(priority)

	// Abrir handle do processo
	handle, err := windows.OpenProcess(PROCESS_SET_INFORMATION, false, uint32(targetProcessID))
	if err != nil {
		return fmt.Errorf("erro ao abrir processo %d: %w", targetProcessID, err)
	}
	defer windows.CloseHandle(handle)

	// Definir prioridade
	err = windows.SetPriorityClass(handle, windowsPriority)
	if err != nil {
		return fmt.Errorf("erro ao definir prioridade do processo %d: %w", targetProcessID, err)
	}

	return nil
}

// GetProcessPriority obtém a prioridade de um processo
func (s *windowsProcessPriorityService) GetProcessPriority(ctx context.Context, processID int) (int, error) {

	handle, err := windows.OpenProcess(PROCESS_QUERY_INFORMATION, false, uint32(processID))
	if err != nil {
		return 0, fmt.Errorf("erro ao abrir processo %d: %w", processID, err)
	}
	defer windows.CloseHandle(handle)

	// Obter classe de prioridade
	priorityClass, err := windows.GetPriorityClass(handle)
	if err != nil {
		return 0, fmt.Errorf("erro ao obter prioridade do processo %d: %w", processID, err)
	}

	// Converter para nosso sistema de prioridade (0-5)
	return s.convertFromWindowsPriority(priorityClass), nil
}

func (s *windowsProcessPriorityService) OptimizeGameProcesses(ctx context.Context) error {
	// Lista de processos comuns de jogos e suas prioridades desejadas
	gameProcesses := map[string]int{
		// Launchers de jogos
		"steam.exe":           3, // Above Normal
		"epicgameslauncher.exe": 3,
		"origin.exe":          3,
		"uplay.exe":           3,
		"battlenet.exe":       3,
		
		// Processos que devem ter prioridade baixa
		"chrome.exe":          1, // Below Normal
		"firefox.exe":         1,
		"discord.exe":         2, // Normal
		"spotify.exe":         1,
		"teams.exe":           1,
		"skype.exe":           1,
		
		// Processos do sistema que podem ser otimizados
		"dwm.exe":             3, // Desktop Window Manager
		"audiodg.exe":         4, // Audio Device Graph
	}

	processes, err := s.GetRunningProcesses(ctx)
	if err != nil {
		return fmt.Errorf("erro ao obter processos: %w", err)
	}

	var errors []string
	for _, proc := range processes {
		if priority, exists := gameProcesses[strings.ToLower(proc.Name)]; exists {
			if err := s.setProcessPriorityByPID(proc.PID, priority); err != nil {
				errors = append(errors, fmt.Sprintf("erro ao otimizar %s (PID %d): %v", proc.Name, proc.PID, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("erros durante otimização: %s", strings.Join(errors, "; "))
	}

	return nil
}

// SetProcessAffinity define a afinidade de CPU de um processo
func (s *windowsProcessPriorityService) SetProcessAffinity(ctx context.Context, processID int, cpuMask uint64) error {
	handle, err := windows.OpenProcess(PROCESS_SET_INFORMATION, false, uint32(processID))
	if err != nil {
		return fmt.Errorf("erro ao abrir processo %d: %w", processID, err)
	}
	defer windows.CloseHandle(handle)

	// Usar SetProcessAffinityMask
	setProcessAffinityMask := s.kernel32.NewProc("SetProcessAffinityMask")
	ret, _, err := setProcessAffinityMask.Call(
		uintptr(handle),
		uintptr(cpuMask),
	)

	if ret == 0 {
		return fmt.Errorf("erro ao definir afinidade do processo %d: %w", processID, err)
	}

	return nil
}

// GetProcessAffinity obtém a afinidade de CPU de um processo
func (s *windowsProcessPriorityService) GetProcessAffinity(ctx context.Context, processID int) (uint64, error) {
	handle, err := windows.OpenProcess(PROCESS_QUERY_INFORMATION, false, uint32(processID))
	if err != nil {
		return 0, fmt.Errorf("erro ao abrir processo %d: %w", processID, err)
	}
	defer windows.CloseHandle(handle)

	// Usar GetProcessAffinityMask
	getProcessAffinityMask := s.kernel32.NewProc("GetProcessAffinityMask")
	var processAffinityMask, systemAffinityMask uintptr
	
	ret, _, err := getProcessAffinityMask.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&processAffinityMask)),
		uintptr(unsafe.Pointer(&systemAffinityMask)),
	)

	if ret == 0 {
		return 0, fmt.Errorf("erro ao obter afinidade do processo %d: %w", processID, err)
	}

	return uint64(processAffinityMask), nil
}

func (s *windowsProcessPriorityService) GetRunningProcesses(ctx context.Context) ([]*entities.ProcessInfo, error) {
	// Criar snapshot dos processos
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar snapshot: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	var processes []*entities.ProcessInfo
	var entry processEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	// Usar Process32First e Process32Next
	process32First := s.kernel32.NewProc("Process32FirstW")
	process32Next := s.kernel32.NewProc("Process32NextW")

	ret, _, _ := process32First.Call(uintptr(snapshot), uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("erro ao obter primeiro processo")
	}

	for {
		// Converter nome do executável de UTF-16 para string
		exeName := windows.UTF16ToString(entry.ExeFile[:])
		
		// Obter caminho completo do executável
		executablePath := s.getProcessExecutablePath(int(entry.ProcessID))

		// Obter informações de memória
		memoryInfo, _ := s.getProcessMemoryInfo(int(entry.ProcessID))

		processInfo := &entities.ProcessInfo{
			PID:            int(entry.ProcessID),
			PPID:           int(entry.ParentProcessID),
			Name:           exeName,
			Path: executablePath,
			Priority:       s.convertFromWindowsPriority(uint32(entry.PriorityClassBase)),
			MemoryUsage:    memoryInfo,
			ThreadCount:    int(entry.Threads),
		}

		processes = append(processes, processInfo)

		// Próximo processo
		ret, _, _ = process32Next.Call(uintptr(snapshot), uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}

	return processes, nil
}

// KillProcess encerra um processo
func (s *windowsProcessPriorityService) KillProcess(ctx context.Context, processID int) error {
	handle, err := windows.OpenProcess(PROCESS_TERMINATE, false, uint32(processID))
	if err != nil {
		return fmt.Errorf("erro ao abrir processo %d: %w", processID, err)
	}
	defer windows.CloseHandle(handle)

	err = windows.TerminateProcess(handle, 1)
	if err != nil {
		return fmt.Errorf("erro ao encerrar processo %d: %w", processID, err)
	}

	return nil
}

// IsProcessRunning verifica se um processo está em execução
func (s *windowsProcessPriorityService) IsProcessRunning(ctx context.Context, processName string) (bool, error) {
	processes, err := s.GetRunningProcesses(ctx)
	if err != nil {
		return false, err
	}

	processName = strings.ToLower(processName)
	for _, proc := range processes {
		if strings.ToLower(proc.Name) == processName {
			return true, nil
		}
	}

	return false, nil
}

// Métodos auxiliares

func (s *windowsProcessPriorityService) convertToWindowsPriority(priority int) uint32 {
	switch priority {
	case 0:
		return IDLE_PRIORITY_CLASS
	case 1:
		return BELOW_NORMAL_PRIORITY_CLASS
	case 2:
		return NORMAL_PRIORITY_CLASS
	case 3:
		return ABOVE_NORMAL_PRIORITY_CLASS
	case 4:
		return HIGH_PRIORITY_CLASS
	case 5:
		return REALTIME_PRIORITY_CLASS
	default:
		return NORMAL_PRIORITY_CLASS
	}
}

func (s *windowsProcessPriorityService) convertFromWindowsPriority(windowsPriority uint32) int {
	switch windowsPriority {
	case IDLE_PRIORITY_CLASS:
		return 0
	case BELOW_NORMAL_PRIORITY_CLASS:
		return 1
	case NORMAL_PRIORITY_CLASS:
		return 2
	case ABOVE_NORMAL_PRIORITY_CLASS:
		return 3
	case HIGH_PRIORITY_CLASS:
		return 4
	case REALTIME_PRIORITY_CLASS:
		return 5
	default:
		return 2
	}
}

func (s *windowsProcessPriorityService) setProcessPriorityByPID(processID int, priority int) error {
	windowsPriority := s.convertToWindowsPriority(priority)

	handle, err := windows.OpenProcess(PROCESS_SET_INFORMATION, false, uint32(processID))
	if err != nil {
		return fmt.Errorf("erro ao abrir processo %d: %w", processID, err)
	}
	defer windows.CloseHandle(handle)

	err = windows.SetPriorityClass(handle, windowsPriority)
	if err != nil {
		return fmt.Errorf("erro ao definir prioridade do processo %d: %w", processID, err)
	}

	return nil
}

func (s *windowsProcessPriorityService) getProcessExecutablePath(processID int) string {
	handle, err := windows.OpenProcess(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ, false, uint32(processID))
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(handle)

	// Buffer para o caminho
	var buffer [windows.MAX_PATH]uint16
	size := uint32(len(buffer))

	// Usar QueryFullProcessImageName
	queryFullProcessImageName := s.kernel32.NewProc("QueryFullProcessImageNameW")
	ret, _, _ := queryFullProcessImageName.Call(
		uintptr(handle),
		uintptr(0), // PROCESS_NAME_WIN32
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret != 0 {
		return windows.UTF16ToString(buffer[:size])
	}

	return ""
}

func (s *windowsProcessPriorityService) getProcessMemoryInfo(processID int) (uint64, error) {
	handle, err := windows.OpenProcess(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ, false, uint32(processID))
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(handle)

	// Estrutura PROCESS_MEMORY_COUNTERS
	type processMemoryCounters struct {
		cb                         uint32
		pageFaultCount             uint32
		peakWorkingSetSize         uintptr
		workingSetSize             uintptr
		quotaPeakPagedPoolUsage    uintptr
		quotaPagedPoolUsage        uintptr
		quotaPeakNonPagedPoolUsage uintptr
		quotaNonPagedPoolUsage     uintptr
		pagefileUsage              uintptr
		peakPagefileUsage          uintptr
	}

	var memCounters processMemoryCounters
	memCounters.cb = uint32(unsafe.Sizeof(memCounters))

	psapi := windows.NewLazySystemDLL("psapi.dll")
	getProcessMemoryInfo := psapi.NewProc("GetProcessMemoryInfo")

	ret, _, _ := getProcessMemoryInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&memCounters)),
		uintptr(memCounters.cb),
	)

	if ret != 0 {
		return uint64(memCounters.workingSetSize), nil
	}

	return 0, fmt.Errorf("erro ao obter informações de memória")
}
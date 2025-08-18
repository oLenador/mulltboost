package outbound

import (
	"context"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type ProcessPriorityService interface {
	// Prioridades de processo
	SetProcessPriority(ctx context.Context, processName string, priority int) error
	GetProcessPriority(ctx context.Context, processID int) (int, error)
	OptimizeGameProcesses(ctx context.Context) error
	
	// Afinidade de CPU
	SetProcessAffinity(ctx context.Context, processID int, cpuMask uint64) error
	GetProcessAffinity(ctx context.Context, processID int) (uint64, error)
	
	// Gerenciamento
	GetRunningProcesses(ctx context.Context) ([]*entities.ProcessInfo, error)
	KillProcess(ctx context.Context, processID int) error
	IsProcessRunning(ctx context.Context, processName string) (bool, error)
}
package windows

import (
	"context"
)

type MemoryManagementService interface {
	// Gerenciamento de memória
	ClearMemoryCache(ctx context.Context) error
	OptimizeMemoryUsage(ctx context.Context) error
	SetVirtualMemorySize(ctx context.Context, sizeGB int) error
	
	// Análise
	GetMemoryInfo(ctx context.Context) (*entities.MemoryInfo, error)
	GetMemoryUsageByProcess(ctx context.Context) (map[string]float64, error)
	FindMemoryLeaks(ctx context.Context) ([]string, error)
	
	// Configuração
	EnableMemoryCompression(ctx context.Context) error
	DisableMemoryCompression(ctx context.Context) error
	SetPageFileLocation(ctx context.Context, drive string) error
}
package windows

import "context"

type CacheManagementService interface {
	// Limpeza de cache
	ClearSystemCache(ctx context.Context) error
	ClearDNSCache(ctx context.Context) error
	ClearTemporaryFiles(ctx context.Context) error
	ClearBrowserCache(ctx context.Context) error
	
	// Análise
	AnalyzeCacheUsage(ctx context.Context) (int64, error)
	GetCacheStatistics(ctx context.Context) (*entities.CacheData, error)
	
	// Configuração
	SetCacheCleanupSchedule(ctx context.Context, schedule string) error
	DisableAutomaticCacheCleanup(ctx context.Context) error
}
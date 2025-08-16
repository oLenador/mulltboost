package windows

import (
	"context"
)

type RegistryService interface {
	// Operações básicas
	ReadRegistryValue(ctx context.Context, keyPath, valueName string) (interface{}, error)
	WriteRegistryValue(ctx context.Context, keyPath, valueName string, value interface{}) error
	DeleteRegistryValue(ctx context.Context, keyPath, valueName string) error
	
	// Chaves
	CreateRegistryKey(ctx context.Context, keyPath string) error
	DeleteRegistryKey(ctx context.Context, keyPath string) error
	KeyExists(ctx context.Context, keyPath string) (bool, error)
	
	// Backup e restore
	BackupRegistryKey(ctx context.Context, keyPath, backupPath string) error
	RestoreRegistryKey(ctx context.Context, keyPath, backupPath string) error
	
	// Otimização
	ApplyPerformanceRegistryTweaks(ctx context.Context) error
	RestoreDefaultRegistrySettings(ctx context.Context) error
}
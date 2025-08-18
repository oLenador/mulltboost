package outbound

import (
	"context"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type NetworkAdapterService interface {
	// Adaptadores de rede
	GetNetworkAdapters(ctx context.Context) ([]*entities.NetworkAdapter, error)
	ResetNetworkStack(ctx context.Context) error
	
	// Configurações
	SetNetworkPriority(ctx context.Context, adapterID string, priority int) error
	DisableNetworkThrottling(ctx context.Context) error
	
	// Status
	GetNetworkStatistics(ctx context.Context, adapterID string) (*entities.NetworkAdapter, error)
	IsNetworkOptimized(ctx context.Context) (bool, error)
}
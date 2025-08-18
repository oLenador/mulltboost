package outbound

import (
	"context"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type TCPOptimizationService interface {
	// Otimização de TCP
	OptimizeTCPForGaming(ctx context.Context) error
	OptimizeTCPForStreaming(ctx context.Context) error
	RestoreDefaultTCPSettings(ctx context.Context) error
	
	// Configurações específicas
	SetTCPWindowSize(ctx context.Context, size int) error
	SetTCPCongestionControl(ctx context.Context, algorithm string) error
	DisableNagleAlgorithm(ctx context.Context) error
	
	// Status
	GetTCPConfiguration(ctx context.Context) (*entities.TCPConfiguration, error)
	IsTCPOptimized(ctx context.Context) (bool, error)
	GetNetworkLatency(ctx context.Context, target string) (int, error)
}
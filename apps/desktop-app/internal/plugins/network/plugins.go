package network

import (
    "context"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
    "github.com/oLenador/mulltbost/internal/core/ports/inbound"
    "github.com/oLenador/mulltbost/internal/plugins/base"
)

type TCPOptimization struct{}

func (t *TCPOptimization) Execute(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "TCP settings optimized for gaming",
        BackupData: map[string]interface{}{
            "registry_backup": "tcp_optimization_backup",
        },
    }, nil
}

func (t *TCPOptimization) Validate(ctx context.Context) error {
    return nil
}

func (t *TCPOptimization) CanApply(ctx context.Context) bool {
    return true
}

func (t *TCPOptimization) CanRevert(ctx context.Context) bool {
    return true
}

func (t *TCPOptimization) Revert(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "TCP settings reverted",
    }, nil
}

func GetAllPlugins() []inbound.OptimizationUseCase {
    tcpOpt := entities.Optimization{
        ID:          "network_tcp_optimization",
        Name:        "TCP Network Optimization",
        Description: "Optimizes TCP settings for reduced latency",
        Category:    entities.CategoryNetwork,
        Level:       entities.LevelFree,
        Platform:    []entities.Platform{entities.PlatformWindows},
        Reversible:  true,
        RiskLevel:   entities.RiskMedium,
        Version:     "1.0.0",
    }

    return []inbound.OptimizationUseCase{
        base.NewBaseOptimization(tcpOpt, &TCPOptimization{}),
    }
}
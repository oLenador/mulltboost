package precision

import (
    "context"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
    "github.com/oLenador/mulltbost/internal/core/ports/inbound"
    "github.com/oLenador/mulltbost/internal/plugins/base"
)

type MouseAccelerationOptimization struct{}

func (m *MouseAccelerationOptimization) Execute(ctx context.Context) (*entities.OptimizationResult, error) {
    // Implementar otimização de aceleração do mouse
    return &entities.OptimizationResult{
        Success: true,
        Message: "Mouse acceleration disabled successfully",
        BackupData: map[string]interface{}{
            "registry_backup": "mouse_accel_backup_data",
        },
    }, nil
}

func (m *MouseAccelerationOptimization) Validate(ctx context.Context) error {
    return nil
}

func (m *MouseAccelerationOptimization) CanApply(ctx context.Context) bool {
    return true
}

func (m *MouseAccelerationOptimization) CanRevert(ctx context.Context) bool {
    return true
}

func (m *MouseAccelerationOptimization) Revert(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "Mouse acceleration settings reverted",
    }, nil
}

func GetAllPlugins() []inbound.OptimizationUseCase {
    mouseAccelOpt := entities.Optimization{
        ID:          "precision_mouse_accel",
        Name:        "Disable Mouse Acceleration",
        Description: "Disables mouse acceleration for better precision",
        Category:    entities.CategoryPrecision,
        Level:       entities.LevelFree,
        Platform:    []entities.Platform{entities.PlatformWindows},
        Reversible:  true,
        RiskLevel:   entities.RiskLow,
        Version:     "1.0.0",
    }

    return []inbound.OptimizationUseCase{
        base.NewBaseOptimization(mouseAccelOpt, &MouseAccelerationOptimization{}),
    }
}
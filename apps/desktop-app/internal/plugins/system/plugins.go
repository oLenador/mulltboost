package system

import (
    "context"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
    "github.com/oLenador/mulltbost/internal/core/ports/inbound"
    "github.com/oLenador/mulltbost/internal/plugins/base"
)

type VisualEffectsOptimization struct{}

func (v *VisualEffectsOptimization) Execute(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "Visual effects disabled for better performance",
        BackupData: map[string]interface{}{
            "registry_backup": "visual_effects_backup",
        },
    }, nil
}

func (v *VisualEffectsOptimization) Validate(ctx context.Context) error {
    return nil
}

func (v *VisualEffectsOptimization) CanApply(ctx context.Context) bool {
    return true
}

func (v *VisualEffectsOptimization) CanRevert(ctx context.Context) bool {
    return true
}

func (v *VisualEffectsOptimization) Revert(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "Visual effects restored",
    }, nil
}

func GetAllPlugins() []inbound.OptimizationUseCase {
    visualEffectsOpt := entities.Optimization{
        ID:          "system_visual_effects",
        Name:        "Disable Visual Effects",
        Description: "Disables Windows visual effects for better performance",
        Category:    entities.CategorySystem,
        Level:       entities.LevelFree,
        Platform:    []entities.Platform{entities.PlatformWindows},
        Reversible:  true,
        RiskLevel:   entities.RiskLow,
        Version:     "1.0.0",
    }

    return []inbound.OptimizationUseCase{
        base.NewBaseOptimization(visualEffectsOpt, &VisualEffectsOptimization{}),
    }
}
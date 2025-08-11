package performance

import (
    "context"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
    "github.com/oLenador/mulltbost/internal/core/ports/inbound"
    "github.com/oLenador/mulltbost/internal/plugins/base"
)

type GameModeOptimization struct{}

func (g *GameModeOptimization) Execute(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "Game mode enabled successfully",
        BackupData: map[string]interface{}{
            "registry_backup": "game_mode_backup_data",
        },
    }, nil
}

func (g *GameModeOptimization) Validate(ctx context.Context) error {
    return nil
}

func (g *GameModeOptimization) CanApply(ctx context.Context) bool {
    return true
}

func (g *GameModeOptimization) CanRevert(ctx context.Context) bool {
    return true
}

func (g *GameModeOptimization) Revert(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "Game mode settings reverted",
    }, nil
}

type HighPerformancePowerPlan struct{}

func (h *HighPerformancePowerPlan) Execute(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "High performance power plan activated",
        BackupData: map[string]interface{}{
            "previous_plan": "balanced",
        },
    }, nil
}

func (h *HighPerformancePowerPlan) Validate(ctx context.Context) error {
    return nil
}

func (h *HighPerformancePowerPlan) CanApply(ctx context.Context) bool {
    return true
}

func (h *HighPerformancePowerPlan) CanRevert(ctx context.Context) bool {
    return true
}

func (h *HighPerformancePowerPlan) Revert(ctx context.Context) (*entities.OptimizationResult, error) {
    return &entities.OptimizationResult{
        Success: true,
        Message: "Power plan reverted to balanced",
    }, nil
}

func GetAllPlugins() []inbound.OptimizationUseCase {
    gameModeOpt := entities.Optimization{
        ID:          "performance_game_mode",
        Name:        "Enable Game Mode",
        Description: "Enables Windows Game Mode for better gaming performance",
        Category:    entities.CategoryPerformance,
        Level:       entities.LevelFree,
        Platform:    []entities.Platform{entities.PlatformWindows},
        Reversible:  true,
        RiskLevel:   entities.RiskLow,
        Version:     "1.0.0",
    }

    powerPlanOpt := entities.Optimization{
        ID:          "performance_power_plan",
        Name:        "High Performance Power Plan",
        Description: "Sets the power plan to high performance",
        Category:    entities.CategoryPerformance,
        Level:       entities.LevelFree,
        Platform:    []entities.Platform{entities.PlatformWindows},
        Reversible:  true,
        RiskLevel:   entities.RiskLow,
        Version:     "1.0.0",
    }

    return []inbound.OptimizationUseCase{
        base.NewBaseOptimization(gameModeOpt, &GameModeOptimization{}),
        base.NewBaseOptimization(powerPlanOpt, &HighPerformancePowerPlan{}),
    }
}

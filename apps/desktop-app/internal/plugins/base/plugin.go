package base

import (
    "context"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
    "github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

type BaseOptimization struct {
    info     entities.Optimization
    executor OptimizationExecutor
}

type OptimizationExecutor interface {
    Execute(ctx context.Context) (*entities.OptimizationResult, error)
    Validate(ctx context.Context) error
    CanApply(ctx context.Context) bool
    CanRevert(ctx context.Context) bool
    Revert(ctx context.Context) (*entities.OptimizationResult, error)
}

func NewBaseOptimization(info entities.Optimization, executor OptimizationExecutor) inbound.OptimizationUseCase {
    return &BaseOptimization{
        info:     info,
        executor: executor,
    }
}

func (b *BaseOptimization) Execute(ctx context.Context) (*entities.OptimizationResult, error) {
    return b.executor.Execute(ctx)
}

func (b *BaseOptimization) Validate(ctx context.Context) error {
    return b.executor.Validate(ctx)
}

func (b *BaseOptimization) CanApply(ctx context.Context) bool {
    return b.executor.CanApply(ctx)
}

func (b *BaseOptimization) CanRevert(ctx context.Context) bool {
    return b.executor.CanRevert(ctx)
}

func (b *BaseOptimization) GetInfo() entities.Optimization {
    return b.info
}

func (b *BaseOptimization) Revert(ctx context.Context) (*entities.OptimizationResult, error) {
    return b.executor.Revert(ctx)
}
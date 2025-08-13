package interfaces

import (
	"context"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type BoosterStrategy interface {
    Execute(ctx context.Context) (*entities.BoosterResult, error)
    Validate(ctx context.Context) error
    CanApply(ctx context.Context) bool
    CanRevert(ctx context.Context) bool
    Revert(ctx context.Context) (*entities.BoosterResult, error)
}

type BoosterFactory interface {
    CreateOptimization(id string) (BoosterStrategy, error)
    GetSupportedOptimizations() []string
}

type BoosterCommand interface {
    Execute(ctx context.Context) (*entities.BoosterResult, error)
    Undo(ctx context.Context) (*entities.BoosterResult, error)
    GetInfo() entities.Booster
}

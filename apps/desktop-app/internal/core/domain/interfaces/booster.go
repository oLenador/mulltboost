package interfaces

import (
	"context"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type BoosterStrategy interface {
    Execute(ctx context.Context) (*entities.BoostApplyResult, error)
    Validate(ctx context.Context) error
    CanApply(ctx context.Context) bool
    CanRevert(ctx context.Context) bool
    Revert(ctx context.Context) (*entities.BoostRevertResult, error)
}

type BoosterFactory interface {
    CreateOptimization(id string) (BoosterStrategy, error)
    GetSupportedOptimizations() []string
}

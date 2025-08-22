package dto

import (
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)


type GetBoosterDto struct {
	ID           string
	Name         string
	Description  string
	Category     entities.BoosterCategory
	Level        entities.BoosterLevel
	Platform     []entities.Platform
	Dependencies []string
	Conflicts    []string
	Reversible   bool
	RiskLevel    entities.RiskLevel
	Version      string
	IsApplied    bool
	AppliedAt    *time.Time
	RevertedAt   *time.Time
	Tags         []string
}


type BoosterDto struct {
	ID           string
	Name         string
	Description  string
	Category     entities.BoosterCategory
	Level        entities.BoosterLevel
	Platform     []entities.Platform
	Dependencies []string
	Conflicts    []string
	Reversible   bool
	RiskLevel    entities.RiskLevel
	Version      string
	Tags         []string
}

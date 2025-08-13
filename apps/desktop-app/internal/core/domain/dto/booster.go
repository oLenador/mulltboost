package dto

import "github.com/oLenador/mulltbost/internal/core/domain/entities"

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

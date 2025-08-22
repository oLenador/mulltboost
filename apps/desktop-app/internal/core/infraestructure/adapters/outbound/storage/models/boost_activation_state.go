package storage

import (
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type BoostActivationState struct {
	ID           string                        `gorm:"primaryKey;type:text"`
	IsApplied    bool                          `gorm:"not null;default:false"`
	AppliedAt    *time.Time                    `gorm:"index"`
	RevertedAt   *time.Time                    `gorm:"index"`
	Version      string                        `gorm:"type:text;not null"`
	Status       entities.BoosterExecutionStatus `gorm:"type:text;not null;index"`
	ErrorMessage string                        `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (BoostActivationState) TableName() string {
	return "boost_activation_states"
}

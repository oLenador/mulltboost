package storage

import (
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type BoostOperation struct {
	ID         string                        `gorm:"primaryKey;type:text"`
	UserID     string                        `gorm:"type:text;index"`
	BoosterID  string                        `gorm:"type:text;index"`
	Version    string                        `gorm:"type:text"`
	AppliedAt  time.Time                     `gorm:"not null;index"`
	RevertedAt time.Time                    `gorm:"index"`
	Status     entities.BoosterExecutionStatus      `gorm:"type:text;index"`
	ErrorMsg   string                        `gorm:"type:text"`
	Type       entities.BoosterOperationType `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (BoostOperation) TableName() string { return "boosts_operations" }

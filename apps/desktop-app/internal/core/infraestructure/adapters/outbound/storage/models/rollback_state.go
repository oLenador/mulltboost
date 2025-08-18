package storage

import (
	"gorm.io/datatypes"
	"time"
)

type BoosterRollbackState struct {
	ID         string            `gorm:"primaryKey;type:text"`
	Applied    bool              `gorm:"not null;default:false"`
	AppliedAt  *time.Time        `gorm:"index"`
	RevertedAt *time.Time        `gorm:"index"`
	Version    string            `gorm:"type:text;not null"`
	BackupData datatypes.JSONMap `gorm:"type:json;not null;default:'{}'"`
	Status     ExecutionStatus   `gorm:"type:text;not null;index"`
	ErrorMsg   string            `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (BoosterRollbackState) TableName() string { return "rollback_states" }

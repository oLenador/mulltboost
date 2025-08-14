package storage

import "time"

type AppliedBoost struct {
	ID         string            `gorm:"primaryKey;type:text"`
	UserID     string            `gorm:"type:text;index"`
	BoosterID  string            `gorm:"type:text;index"`
	Version    string            `gorm:"type:text"`
	AppliedAt  time.Time         `gorm:"not null;index"`
	RevertedAt *time.Time        `gorm:"index"`
	Status     ExecutionStatus  `gorm:"type:text;index"`
	ErrorMsg   string            `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (AppliedBoost) TableName() string { return "applied_boosts" }

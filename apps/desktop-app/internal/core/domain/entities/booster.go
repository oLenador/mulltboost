package entities

import (
	"time"
)

type BoosterCategory string

const (
	CategoryConnection BoosterCategory = "connection"
	CategoryFlusher    BoosterCategory = "flusher"
	CategoryFPSBooster BoosterCategory = "fps-booster"
	CategoryGames      BoosterCategory = "games"
	CategoryPrecision  BoosterCategory = "precision"
)

type BoosterOperationType string

const (
	RevertOperationType BoosterLevel = "revert"
	ApplyOperationType  BoosterLevel = "apply"
)

type BoosterLevel string

const (
	LevelFree    BoosterLevel = "free"
	LevelPremium BoosterLevel = "premium"
)

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type Platform string

const (
	PlatformWindows Platform = "windows"
	PlatformLinux   Platform = "linux"
)

type ExecutionStatus string

const (
	BoosterStatusNotApplied ExecutionStatus = "not_applied"
	BoosterStatusApplied    ExecutionStatus = "applied"
	BoosterStatusApplying   ExecutionStatus = "applying"
	BoosterStatusFailed     ExecutionStatus = "failed"
	BoosterStatusReverting  ExecutionStatus = "reverting"
	BoosterStatusReverted   ExecutionStatus = "reverted"
)

type BoosterOpeCallResult struct {
	OperationID string
	SubmittedAt time.Time
	Status      OperationStatus
}

type OperationStatus string

const (
	OperationStatusPending    OperationStatus = "pending"
	OperationStatusProcessing OperationStatus = "processing"
	OperationStatusFailed     OperationStatus = "failed"
)

type BoostApplyResult struct {
	Success    bool
	Message    string
	BackupData map[string]interface{}
	Error      error
}

type BoostRevertResult struct {
	Success    bool
	Message    string
	Error      error
}


type Booster struct {
	ID             string
	NameKey        string
	DescriptionKey string
	Category       BoosterCategory
	Level          BoosterLevel
	Platform       []Platform
	Dependencies   []string
	Conflicts      []string
	Reversible     bool
	RiskLevel      RiskLevel
	Version        string
	Tags           []string
}

type BackupData map[string]interface{}
type BoosterRollbackState struct {
	ID         string
	Applied    bool
	AppliedAt  *time.Time
	RevertedAt *time.Time
	Version    string
	BackupData BackupData
	Status     ExecutionStatus
	ErrorMsg   string
}

type BoostOperation struct {
	ID         string
	UserID     string
	BoosterID  string
	Version    string
	Type       BoosterOperationType
	AppliedAt  time.Time
	RevertedAt *time.Time
	Status     ExecutionStatus
	ErrorMsg   string
}

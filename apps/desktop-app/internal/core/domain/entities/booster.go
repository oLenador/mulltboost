package entities

import (
	"context"
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

type BaseStatus string
type BoosterExecutionStatus BaseStatus
type OperationStatus BaseStatus

// Booster Execution Status - Estados do booster
const (
	ExecutionNotApplied BoosterExecutionStatus = "not_applied"
	ExecutionApplied    BoosterExecutionStatus = "applied"
	ExecutionApplying   BoosterExecutionStatus = "applying"
	ExecutionFailed     BoosterExecutionStatus = "failed"
	ExecutionReverting  BoosterExecutionStatus = "reverting"
	ExecutionPending  BoosterExecutionStatus = "pending"
	ExecutionReverted  BoosterExecutionStatus = "reverted"

	ExecutionInactive  BoosterExecutionStatus = "Inactive"
)

const (
	StatusActive   = ExecutionApplied
	StatusInactive = ExecutionInactive
	StatusObsolete BoosterExecutionStatus = "obsolete"
)

// Operation Status - Estados das operações
const (
	OperationPending    OperationStatus = "pending"
	OperationProcessing OperationStatus = "processing"
	OperationCompleted  OperationStatus = "completed"
	OperationFailed     OperationStatus = "failed"
	OperationCancelled  OperationStatus = "cancelled"
)

type BoosterOperationType string

const (
	RevertOperationType BoosterOperationType = "revert"
	ApplyOperationType  BoosterOperationType = "apply"
)

type BoosterOpeCallResult struct {
	OperationID string
	SubmittedAt time.Time
	Status      OperationStatus
}
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
	Status     BoosterExecutionStatus
	ErrorMsg   string
}

type BoostOperation struct {
	ID         string
	BoosterID  string
	Type       BoosterOperationType
	AppliedAt  time.Time
	RevertedAt time.Time
	ErrorMsg   string
}


type QueueState struct {
	Items []QueueItem
	QueueSize int
	TotalProcessed int
	InProgress int
}

type QueueItem struct {
	ID          string
	BoosterID   string
	Operation   BoosterOperationType
	OperationID string
	SubmittedAt time.Time
	Context     context.Context
	Cancel      context.CancelFunc
}
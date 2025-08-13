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
	StatusNotApplied ExecutionStatus = "not_applied"
	StatusApplied    ExecutionStatus = "applied"
	StatusFailed     ExecutionStatus = "failed"
	StatusReverting  ExecutionStatus = "reverting"
	StatusReverted   ExecutionStatus = "reverted"
)

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

type BoosterState struct {
	ID         string
	Applied    bool
	AppliedAt  *time.Time
	RevertedAt *time.Time
	Version    string
	BackupData map[string]interface{}
	Status     ExecutionStatus
	ErrorMsg   string
}

type BoosterResult struct {
	Success    bool
	Message    string
	BackupData map[string]interface{}
	Error      error
}

type BatchResult struct {
	TotalCount    int
	SuccessCount  int
	FailedCount   int
	Results       map[string]BoosterResult
	OverallStatus string
}

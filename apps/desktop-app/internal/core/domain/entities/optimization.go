package entities

import (
    "time"
)

type OptimizationCategory string

const (
    CategoryPrecision    OptimizationCategory = "precision"
    CategoryPerformance  OptimizationCategory = "performance"  
    CategoryNetwork      OptimizationCategory = "network"
    CategorySystem       OptimizationCategory = "system"
)

type OptimizationLevel string

const (
    LevelFree    OptimizationLevel = "free"
    LevelPremium OptimizationLevel = "premium"
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

type Optimization struct {
    ID          string
    Name        string
    Description string
    Category    OptimizationCategory
    Level       OptimizationLevel
    Platform    []Platform
    Dependencies []string
    Conflicts   []string
    Reversible  bool
    RiskLevel   RiskLevel
    Version     string
}

type OptimizationState struct {
    ID         string
    Applied    bool
    AppliedAt  *time.Time
    RevertedAt *time.Time
    Version    string
    BackupData map[string]interface{}
    Status     ExecutionStatus
    ErrorMsg   string
}

type OptimizationResult struct {
    Success     bool
    Message     string
    BackupData  map[string]interface{}
    Error       error
}

type BatchResult struct {
    TotalCount    int
    SuccessCount  int
    FailedCount   int
    Results       map[string]OptimizationResult
    OverallStatus string
}
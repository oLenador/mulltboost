package storage

// ExecutionStatus armazenado como string
type ExecutionStatus string

const (
	StatusNotApplied ExecutionStatus = "not_applied"
	StatusApplied    ExecutionStatus = "applied"
	StatusFailed     ExecutionStatus = "failed"
	StatusReverting  ExecutionStatus = "reverting"
	StatusReverted   ExecutionStatus = "reverted"
)
package entities

import "time"

type InitResult struct {
	OperationID string
	SubmittedAt time.Time
	Success     bool
	Status      OperationStatus
	Message     string
	Error       error
}

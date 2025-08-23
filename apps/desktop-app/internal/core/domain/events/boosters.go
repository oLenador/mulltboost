package events

import (
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

const (
	BoosterOperationResult = "booster_operation_result"
)

type BoosterEvent struct {
	EventType     entities.EventStatus
	Timestamp     time.Time
	OperationType entities.BoosterOperationType
	OperationID   string
	BoosterID     string
	Status        entities.BoosterExecutionStatus
	EndAt         time.Time
	Sequency      int
	IdempotencyID string
	Error         error
	QueueSize     int
}

type BoosterBatchProgressEvent struct {
	EventType        entities.EventStatus
	Timestamp        time.Time
	BatchID          string
	OperationType    entities.BoosterOperationType
	TotalCount       int
	QueuedCount      int
	ValidationErrors map[string]error
	QueueSize        int
}

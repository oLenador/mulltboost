package events

import (
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

const (
	BoosterOperationResult = "booster_operation_result"
)

type BoosterEvent struct {
	EventType     string                        `json:"eventType"`
	Timestamp     time.Time                     `json:"timestamp"`
	OperationType entities.BoosterOperationType `json:"operationType"`
	OperationID   string                        `json:"operationId"`
	BoosterID     string                        `json:"boosterId"`
	Status        string                        `json:"status"`
	EndAt         time.Time                     `json:"appliedAt"`
}

type BoosterBatchProgressEvent struct {
	BatchID   string            `json:"batchId"`
	Total     int               `json:"total"`
	Completed int               `json:"completed"`
	Failed    int               `json:"failed"`
	Details   map[string]string `json:"details"`
}

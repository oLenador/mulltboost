package booster

import (
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/wailsapp/wails/v3/pkg/application"
)
type BoosterEventEmitter struct {
	builder *EventBuilder
}

func NewBoosterEventEmitter(eventManager *application.EventManager) *BoosterEventEmitter {
	return &BoosterEventEmitter{
		builder: NewEventBuilder(eventManager),
	}
}

func (e *BoosterEventEmitter) EmitProcessing(boosterID string, operationID string, operation entities.BoosterOperationType) {
	e.builder.EmitProcessing(boosterID, operationID, operation)
}

func (e *BoosterEventEmitter) EmitSuccess(boosterID string, operationID string, operation entities.BoosterOperationType, message string) {
	e.builder.EmitSuccess(boosterID, operationID, operation, message)
}

func (e *BoosterEventEmitter) EmitError(boosterID string, operationID string, operation entities.BoosterOperationType, err error) {
	e.builder.EmitError(boosterID, operationID, operation, err)
}

func (e *BoosterEventEmitter) EmitFailed(boosterID string, operationID string, operation entities.BoosterOperationType, message string) {
	e.builder.EmitFailed(boosterID, operationID, operation, message)
}

// Métodos adicionais para uso no Service
func (e *BoosterEventEmitter) EmitQueued(boosterID string, operationID string, operation entities.BoosterOperationType, queueSize int) {
	e.builder.EmitQueued(boosterID, operationID, operation, queueSize)
}

func (e *BoosterEventEmitter) EmitBatchQueued(batchID string, operation entities.BoosterOperationType, totalCount, queuedCount int, validationErrors map[string]error, queueSize int) {
	e.builder.EmitBatchQueued(batchID, operation, totalCount, queuedCount, validationErrors, queueSize)
}

func (e *BoosterEventEmitter) EmitCancelled(boosterID string, queueSize int) {
	e.builder.EmitCancelled(boosterID, queueSize)
}
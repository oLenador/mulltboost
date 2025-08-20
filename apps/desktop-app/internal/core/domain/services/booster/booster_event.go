package booster

import (
	"fmt"
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/events"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// EventBuilder centraliza a criação e envio de eventos
type EventBuilder struct {
	eventManager *application.EventManager
}

func NewEventBuilder(eventManager *application.EventManager) *EventBuilder {
	return &EventBuilder{
		eventManager: eventManager,
	}
}

// createBoosterEvent cria um evento padronizado
func (eb *EventBuilder) createBoosterEvent(
	eventType entities.EventStatus,
	boosterID string,
	operationID string,
	operation entities.BoosterOperationType,
	status entities.BoosterExecutionStatus,
	errorMsg error,
) *application.CustomEvent {
	return &application.CustomEvent{
		Name: fmt.Sprintf("booster.%s", string(eventType)),
		Data: events.BoosterEvent{
			EventType:     eventType,
			Timestamp:     time.Now(),
			OperationType: operation,
			OperationID:   operationID,
			BoosterID:     boosterID,
			Status:        status,
			Error:         errorMsg,
		},
		Sender: "booster-service",
	}
}

// EmitProcessing emite evento de processamento
func (eb *EventBuilder) EmitProcessing(boosterID, operationID string, operation entities.BoosterOperationType) {
	event := eb.createBoosterEvent(
		entities.EventProcessing,
		boosterID,
		operationID,
		operation,
		entities.ExecutionApplying,
		nil,
	)
	eb.eventManager.EmitEvent(event)
}

// EmitSuccess emite evento de sucesso
func (eb *EventBuilder) EmitSuccess(boosterID, operationID string, operation entities.BoosterOperationType, message string) {
	var status entities.BoosterExecutionStatus

	event := eb.createBoosterEvent(
		entities.EventSuccess,
		boosterID,
		operationID,
		operation,
		status,
		nil,
	)
	eb.eventManager.EmitEvent(event)
}

// EmitError emite evento de erro
func (eb *EventBuilder) EmitError(boosterID, operationID string, operation entities.BoosterOperationType, err error) {
	event := eb.createBoosterEvent(
		entities.EventError,
		boosterID,
		operationID,
		operation,
		entities.ExecutionFailed,
		err,
	)
	eb.eventManager.EmitEvent(event)
}

// EmitFailed emite evento de falha
func (eb *EventBuilder) EmitFailed(boosterID, operationID string, operation entities.BoosterOperationType, message string) {
	event := eb.createBoosterEvent(
		entities.EventFailed,
		boosterID,
		operationID,
		operation,
		entities.ExecutionFailed,
		nil,
	)
	eb.eventManager.EmitEvent(event)
}

// EmitQueued emite evento de enfileiramento
func (eb *EventBuilder) EmitQueued(boosterID, operationID string, operation entities.BoosterOperationType, queueSize int) {
	event := &application.CustomEvent{
		Name: "booster.queued",
		Data: events.BoosterEvent{
			EventType:     entities.EventQueued,
			Timestamp:     time.Now(),
			OperationType: operation,
			OperationID:   operationID,
			BoosterID:     boosterID,
			Status:        entities.ExecutionNotApplied,
			QueueSize:     queueSize,
		},
		Sender: "booster-service",
	}
	eb.eventManager.EmitEvent(event)
}

func (eb *EventBuilder) EmitBatchQueued(
	batchID string,
	operation entities.BoosterOperationType,
	totalCount, queuedCount int,
	validationErrors map[string]error,
	queueSize int,
) {


	event := &application.CustomEvent{
		Name: "booster.batch_queued",
		Data: events.BoosterBatchProgressEvent{
			EventType:         entities.EventQueued,
			Timestamp:         time.Now(),
			BatchID:           batchID,
			OperationType:     operation,
			TotalCount:        totalCount,
			QueuedCount:       queuedCount,
			ValidationErrors:  validationErrors,
			QueueSize:         queueSize,
		},
		Sender: "booster-service",
	}
	eb.eventManager.EmitEvent(event)
}

// EmitCancelled emite evento de cancelamento
func (eb *EventBuilder) EmitCancelled(boosterID string, queueSize int) {
	event := &application.CustomEvent{
		Name: "booster.cancelled",
		Data: events.BoosterEvent{
			EventType: entities.EventCancelled,
			Timestamp: time.Now(),
			BoosterID: boosterID,
			Status:    entities.ExecutionNotApplied,
			QueueSize: queueSize,
		},
		Sender: "booster-service",
	}
	eb.eventManager.EmitEvent(event)
}

// EventEmitter implementa a interface usando EventBuilder
type BoosterEventEmitter struct {
	builder *EventBuilder
}

func NewBoosterEventEmitter(eventManager *application.EventManager) *BoosterEventEmitter {
	return &BoosterEventEmitter{
		builder: NewEventBuilder(eventManager),
	}
}

func (e *BoosterEventEmitter) EmitProcessing(boosterID, operationID string, operation entities.BoosterOperationType) {
	e.builder.EmitProcessing(boosterID, operationID, operation)
}

func (e *BoosterEventEmitter) EmitSuccess(boosterID, operationID string, operation entities.BoosterOperationType, message string) {
	e.builder.EmitSuccess(boosterID, operationID, operation, message)
}

func (e *BoosterEventEmitter) EmitError(boosterID, operationID string, operation entities.BoosterOperationType, err error) {
	e.builder.EmitError(boosterID, operationID, operation, err)
}

func (e *BoosterEventEmitter) EmitFailed(boosterID, operationID string, operation entities.BoosterOperationType, message string) {
	e.builder.EmitFailed(boosterID, operationID, operation, message)
}

// Métodos adicionais para uso no Service
func (e *BoosterEventEmitter) EmitQueued(boosterID, operationID string, operation entities.BoosterOperationType, queueSize int) {
	e.builder.EmitQueued(boosterID, operationID, operation, queueSize)
}

func (e *BoosterEventEmitter) EmitBatchQueued(batchID string, operation entities.BoosterOperationType, totalCount, queuedCount int, validationErrors map[string]error, queueSize int) {
	e.builder.EmitBatchQueued(batchID, operation, totalCount, queuedCount, validationErrors, queueSize)
}

func (e *BoosterEventEmitter) EmitCancelled(boosterID string, queueSize int) {
	e.builder.EmitCancelled(boosterID, queueSize)
}
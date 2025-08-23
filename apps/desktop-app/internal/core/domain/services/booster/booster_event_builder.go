package booster

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/events"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type EventBuilder struct {
	eventManager     *application.EventManager
	sequenceManager  *SequenceManager
	idempotencyCache map[string]*application.CustomEvent
	cacheMutex       sync.RWMutex
	logger           *logger.CustomLogger
}

// EventDataBuilder é o builder unificado para criação de eventos
type EventDataBuilder struct {
	eventBuilder *EventBuilder
	eventType    entities.EventStatus
	boosterID    string
	operationID  string
	operation    entities.BoosterOperationType
	status       entities.BoosterExecutionStatus
	errorMsg     error
	queueSize    int
	batchID      string
	totalCount   int
	queuedCount  int
	validationErrors map[string]error
}

func NewEventBuilder(eventManager *application.EventManager) *EventBuilder {
	return &EventBuilder{
		eventManager:     eventManager,
		sequenceManager:  NewSequenceManager(),
		idempotencyCache: make(map[string]*application.CustomEvent),
		logger:           logger.NewCustomLogger("[EventBuilder]"),
	}
}

// NewEventDataBuilder cria um novo builder de dados de evento
func (eb *EventBuilder) NewEventDataBuilder() *EventDataBuilder {
	return &EventDataBuilder{
		eventBuilder: eb,
	}
}

// Métodos de configuração do builder
func (edb *EventDataBuilder) WithEventType(eventType entities.EventStatus) *EventDataBuilder {
	edb.eventType = eventType
	return edb
}

func (edb *EventDataBuilder) WithBoosterID(boosterID string) *EventDataBuilder {
	edb.boosterID = boosterID
	return edb
}

func (edb *EventDataBuilder) WithOperationID(operationID string) *EventDataBuilder {
	edb.operationID = operationID
	return edb
}

func (edb *EventDataBuilder) WithOperation(operation entities.BoosterOperationType) *EventDataBuilder {
	edb.operation = operation
	return edb
}

func (edb *EventDataBuilder) WithStatus(status entities.BoosterExecutionStatus) *EventDataBuilder {
	edb.status = status
	return edb
}

func (edb *EventDataBuilder) WithError(err error) *EventDataBuilder {
	edb.errorMsg = err
	return edb
}

func (edb *EventDataBuilder) WithQueueSize(queueSize int) *EventDataBuilder {
	edb.queueSize = queueSize
	return edb
}

func (edb *EventDataBuilder) WithBatchID(batchID string) *EventDataBuilder {
	edb.batchID = batchID
	return edb
}

func (edb *EventDataBuilder) WithTotalCount(totalCount int) *EventDataBuilder {
	edb.totalCount = totalCount
	return edb
}

func (edb *EventDataBuilder) WithQueuedCount(queuedCount int) *EventDataBuilder {
	edb.queuedCount = queuedCount
	return edb
}

func (edb *EventDataBuilder) WithValidationErrors(validationErrors map[string]error) *EventDataBuilder {
	edb.validationErrors = validationErrors
	return edb
}

// Build constrói o evento baseado no tipo
func (edb *EventDataBuilder) Build() *application.CustomEvent {
	switch edb.eventType {
	case entities.EventBatchQueued:
		return edb.buildBatchEvent()
	case entities.EventQueued:
		return edb.buildQueuedEvent()
	case entities.EventCancelled:
		return edb.buildCancelledEvent()
	default:
		return edb.buildBoosterEvent()
	}
}

// buildBoosterEvent usa o método createBoosterEvent existente
func (edb *EventDataBuilder) buildBoosterEvent() *application.CustomEvent {
	return edb.eventBuilder.createBoosterEvent(
		edb.eventType,
		edb.boosterID,
		edb.operationID,
		edb.operation,
		edb.status,
		edb.errorMsg,
	)
}

// buildQueuedEvent cria evento de fila usando createBoosterEvent como base
func (edb *EventDataBuilder) buildQueuedEvent() *application.CustomEvent {
	edb.eventBuilder.logger.InfoFields("Creating Queued Event", logger.Fields{
		"boosterID":   edb.boosterID,
		"operationID": edb.operationID,
		"operation":   edb.operation,
		"queueSize":   edb.queueSize,
	})

	// Usa createBoosterEvent como base e ajusta campos específicos
	event := edb.eventBuilder.createBoosterEvent(
		entities.EventQueued,
		edb.boosterID,
		edb.operationID,
		edb.operation,
		entities.ExecutionNotApplied,
		nil,
	)

	// Ajusta campos específicos do evento de fila
	if boosterEvent, ok := event.Data.(events.BoosterEvent); ok {
		boosterEvent.QueueSize = edb.queueSize
		event.Data = boosterEvent
	}

	edb.eventBuilder.logger.InfoFields("Queued EventToReturn", logger.Fields{
		"event": event,
	})

	return event
}

// buildBatchEvent cria evento de lote
func (edb *EventDataBuilder) buildBatchEvent() *application.CustomEvent {
	edb.eventBuilder.logger.InfoFields("Creating Batch Queued Event", logger.Fields{
		"batchID":          edb.batchID,
		"operation":        edb.operation,
		"totalCount":       edb.totalCount,
		"queuedCount":      edb.queuedCount,
		"validationErrors": edb.validationErrors,
		"queueSize":        edb.queueSize,
	})

	event := &application.CustomEvent{
		Name: string(entities.EventBatchQueued),
		Data: events.BoosterBatchProgressEvent{
			EventType:        entities.EventQueued,
			Timestamp:        time.Now(),
			BatchID:          edb.batchID,
			OperationType:    edb.operation,
			TotalCount:       edb.totalCount,
			QueuedCount:      edb.queuedCount,
			ValidationErrors: edb.validationErrors,
			QueueSize:        edb.queueSize,
		},
		Sender: "booster-service",
	}

	edb.eventBuilder.logger.InfoFields("Batch Queued EventToReturn", logger.Fields{
		"event": event,
	})

	return event
}

// buildCancelledEvent cria evento de cancelamento usando createBoosterEvent
func (edb *EventDataBuilder) buildCancelledEvent() *application.CustomEvent {
	edb.eventBuilder.logger.InfoFields("Creating Cancelled Event", logger.Fields{
		"boosterID": edb.boosterID,
		"queueSize": edb.queueSize,
	})

	// Usa createBoosterEvent como base
	event := edb.eventBuilder.createBoosterEvent(
		entities.EventCancelled,
		edb.boosterID,
		"", // operationID vazio para cancelamento
		entities.BoosterOperationType(""), // operation vazio para cancelamento
		entities.ExecutionNotApplied,
		nil,
	)

	// Ajusta campos específicos do evento de cancelamento
	if boosterEvent, ok := event.Data.(events.BoosterEvent); ok {
		boosterEvent.QueueSize = edb.queueSize
		event.Data = boosterEvent
	}

	edb.eventBuilder.logger.InfoFields("Cancelled EventToReturn", logger.Fields{
		"event": event,
	})

	return event
}

// Método original mantido para compatibilidade, mas agora usa o builder interno
func (eb *EventBuilder) createBoosterEvent(
	eventType entities.EventStatus,
	boosterID string,
	operationID string,
	operation entities.BoosterOperationType,
	status entities.BoosterExecutionStatus,
	errorMsg error,
) *application.CustomEvent {
	eb.logger.InfoFields("Creating Booster Event", logger.Fields{
		"eventType":   eventType,
		"boosterID":   boosterID,
		"operationID": operationID,
		"operation":   operation,
		"status":      status,
		"errorMsg":    errorMsg,
	})

	idempotencyKey := fmt.Sprintf("%s:%s:%s", boosterID, operationID, eventType)

	if eb.sequenceManager.IsIdempotent(boosterID, operationID) {
		eb.cacheMutex.RLock()
		if cachedEvent, exists := eb.idempotencyCache[idempotencyKey]; exists {
			eb.logger.InfoFields("Cached event", logger.Fields{
				"event": cachedEvent,
			})
			eb.cacheMutex.RUnlock()
			return cachedEvent
		}
		eb.cacheMutex.RUnlock()
	}

	sequence := eb.sequenceManager.GetNextSequence(boosterID, operationID)
	idempotencyId := uuid.New().String()

	event := &application.CustomEvent{
		Name: string(eventType),
		Data: events.BoosterEvent{
			EventType:     eventType,
			Timestamp:     time.Now(),
			OperationType: operation,
			OperationID:   operationID,
			BoosterID:     boosterID,
			IdempotencyID: idempotencyId,
			Sequency:      int(sequence),
			Status:        status,
			Error:         errorMsg,
		},
		Sender: "booster-service",
	}

	eb.logger.InfoFields("EventToReturn", logger.Fields{
		"event": event,
	})

	eb.cacheMutex.Lock()
	eb.idempotencyCache[idempotencyKey] = event
	eb.cacheMutex.Unlock()

	return event
}

// Métodos de emissão refatorados usando o builder
func (eb *EventBuilder) EmitProcessing(boosterID, operationID string, operation entities.BoosterOperationType) {
	event := eb.NewEventDataBuilder().
		WithEventType(entities.EventProcessing).
		WithBoosterID(boosterID).
		WithOperationID(operationID).
		WithOperation(operation).
		WithStatus(entities.ExecutionApplying).
		Build()
	
	eb.eventManager.EmitEvent(event)
}

func (eb *EventBuilder) EmitSuccess(boosterID, operationID string, operation entities.BoosterOperationType, message string) {
	var status entities.BoosterExecutionStatus
	if operation == entities.ApplyOperationType {
		status = entities.ExecutionApplied
	} else {
		status = entities.ExecutionReverted
	}

	event := eb.NewEventDataBuilder().
		WithEventType(entities.EventSuccess).
		WithBoosterID(boosterID).
		WithOperationID(operationID).
		WithOperation(operation).
		WithStatus(status).
		Build()
	
	eb.eventManager.EmitEvent(event)
	eb.markOperationCompleted(boosterID, operationID)
}

func (eb *EventBuilder) EmitError(boosterID, operationID string, operation entities.BoosterOperationType, err error) {
	event := eb.NewEventDataBuilder().
		WithEventType(entities.EventError).
		WithBoosterID(boosterID).
		WithOperationID(operationID).
		WithOperation(operation).
		WithStatus(entities.ExecutionFailed).
		WithError(err).
		Build()
	
	eb.eventManager.EmitEvent(event)
	eb.markOperationCompleted(boosterID, operationID)
}

func (eb *EventBuilder) EmitFailed(boosterID, operationID string, operation entities.BoosterOperationType, message string) {
	event := eb.NewEventDataBuilder().
		WithEventType(entities.EventFailed).
		WithBoosterID(boosterID).
		WithOperationID(operationID).
		WithOperation(operation).
		WithStatus(entities.ExecutionFailed).
		Build()
	
	eb.eventManager.EmitEvent(event)
	eb.markOperationCompleted(boosterID, operationID)
}

func (eb *EventBuilder) EmitQueued(boosterID, operationID string, operation entities.BoosterOperationType, queueSize int) {
	event := eb.NewEventDataBuilder().
		WithEventType(entities.EventQueued).
		WithBoosterID(boosterID).
		WithOperationID(operationID).
		WithOperation(operation).
		WithQueueSize(queueSize).
		Build()
	
	eb.eventManager.EmitEvent(event)
}

func (eb *EventBuilder) EmitBatchQueued(
	batchID string,
	operation entities.BoosterOperationType,
	totalCount, queuedCount int,
	validationErrors map[string]error,
	queueSize int,
) {
	event := eb.NewEventDataBuilder().
		WithEventType(entities.EventBatchQueued).
		WithBatchID(batchID).
		WithOperation(operation).
		WithTotalCount(totalCount).
		WithQueuedCount(queuedCount).
		WithValidationErrors(validationErrors).
		WithQueueSize(queueSize).
		Build()
	
	eb.eventManager.EmitEvent(event)
}

func (eb *EventBuilder) EmitCancelled(boosterID string, queueSize int) {
	event := eb.NewEventDataBuilder().
		WithEventType(entities.EventCancelled).
		WithBoosterID(boosterID).
		WithQueueSize(queueSize).
		Build()
	
	eb.eventManager.EmitEvent(event)
}

func (eb *EventBuilder) markOperationCompleted(boosterID, operationID string) {
	eb.sequenceManager.MarkCompleted(boosterID)

	// Limpar cache de idempotência para todas as chaves relacionadas à operação
	eb.cacheMutex.Lock()
	defer eb.cacheMutex.Unlock()
	
	keysToDelete := make([]string, 0)
	for key := range eb.idempotencyCache {
		if fmt.Sprintf("%s:%s", boosterID, operationID) == key[:len(fmt.Sprintf("%s:%s", boosterID, operationID))] {
			keysToDelete = append(keysToDelete, key)
		}
	}
	
	for _, key := range keysToDelete {
		delete(eb.idempotencyCache, key)
	}
}
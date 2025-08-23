package booster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/wailsapp/wails/lib/logger"
)

type Processor interface {
	ProcessApply(ctx context.Context, boosterID string) (*entities.BoostApplyResult, error)
	ProcessRevert(ctx context.Context, boosterID string) (*entities.BoostRevertResult, error)
	ValidateBoosterOperation(ctx context.Context, boosterID string, operation entities.BoosterOperationType) error
}

type EventEmitter interface {
	EmitProcessing(boosterID, operationID string, operation entities.BoosterOperationType) 
	EmitSuccess(boosterID, operationID string, operation entities.BoosterOperationType, message string) 
	EmitError(boosterID, operationID string, operation entities.BoosterOperationType, err error) 
	EmitFailed(boosterID, operationID string, operation entities.BoosterOperationType, message string) 
	EmitQueued(boosterID, operationID string, operation entities.BoosterOperationType, queueSize int) 
	EmitBatchQueued(batchID string, operation entities.BoosterOperationType, totalCount, queuedCount int, validationErrors map[string]error, queueSize int) 
	EmitCancelled(boosterID string, queueSize int) 
}

type HistoryRecorder interface {
	RecordOperation(item entities.QueueItem, result *entities.BoostOperation, err error) error
}

// Pool gerencia um pool de workers
type Pool struct {
	processor       Processor
	eventEmitter    EventEmitter
	historyRecorder HistoryRecorder
	queueManager    *Manager
	workerCount     int
	wg              sync.WaitGroup
	stopOnce        sync.Once
}

// NewPool cria um novo pool de workers
func NewPool(
	workerCount int,
	processor Processor,
	eventEmitter EventEmitter,
	historyRecorder HistoryRecorder,
	queueManager *Manager,
) *Pool {
	return &Pool{
		processor:       processor,
		eventEmitter:    eventEmitter,
		historyRecorder: historyRecorder,
		queueManager:    queueManager,
		workerCount:     workerCount,
	}
}

// Start inicia todos os workers
func (p *Pool) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker processa itens da queue
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	workCh := p.queueManager.GetWorkChannel()
	stopCh := p.queueManager.GetStopChannel()

	for {
		select {
		case <-stopCh:
			return
		case item := <-workCh:
			p.processItem(item)
		}
	}
}

func (p *Pool) processItem(item entities.QueueItem) {
	logger.NewCustomLogger("ProcessItem").DebugFields(
		"Listando Listando i",
		logger.Fields{
			"boosters": item,
		},
	)
	p.queueManager.Remove(item.BoosterID)
	p.eventEmitter.EmitProcessing(item.BoosterID, item.OperationID, item.Operation)

	// Validar operação
	if err := p.validateOperation(item); err != nil {
		logger.NewCustomLogger("ProcessItem").DebugFields(
			"Listando Erro",
			logger.Fields{
				"boosters": err,
			},
		)
		p.eventEmitter.EmitError(item.BoosterID, item.OperationID, item.Operation, err)
		return
	}

	// Processar operação
	op, err := p.executeOperation(item)
	logger.NewCustomLogger("ProcessItem").ErrorFields(
		"Listando Erro",
		logger.Fields{
			"boosters": err,
		},
	)

	// Emitir eventos e registrar histórico
	p.handleResult(item, op, err)
}

func (p *Pool) handleResult(item entities.QueueItem, op *entities.BoostOperation, err error) {
	if err != nil {
		p.handleError(item, op, err)
	} else {
		p.handleSuccess(item, op)
	}


	p.historyRecorder.RecordOperation(item, op, err)
}
func (p *Pool) handleError(item entities.QueueItem, op *entities.BoostOperation, err error) {
	p.eventEmitter.EmitError(item.BoosterID, item.OperationID, item.Operation, err)
}


func (p *Pool) handleSuccess(item entities.QueueItem, op *entities.BoostOperation) {
	errorMsg := ""
	if op != nil {
		errorMsg = op.ErrorMsg
	}
	p.eventEmitter.EmitSuccess(item.BoosterID, item.OperationID, item.Operation, errorMsg)
}

// validateOperation valida se a operação pode ser executada
func (p *Pool) validateOperation(item entities.QueueItem) error {
	operationType := entities.BoosterOperationType(item.Operation)

	switch operationType {
	case entities.BoosterOperationType(entities.ApplyOperationType),
		entities.BoosterOperationType(entities.RevertOperationType):
		return p.processor.ValidateBoosterOperation(item.Context, item.BoosterID, operationType)
	default:
		return fmt.Errorf("invalid operation type: %v", item.Operation)
	}
}

// executeOperation executa a operação apropriada
func (p *Pool) executeOperation(item entities.QueueItem) (*entities.BoostOperation, error) {
	switch entities.BoosterOperationType(item.Operation) {
	case entities.BoosterOperationType(entities.ApplyOperationType):
		return p.processApplyOperation(item)
	case entities.BoosterOperationType(entities.RevertOperationType):
		return p.processRevertOperation(item)
	default:
		return nil, fmt.Errorf("unsupported operation type: %v", item.Operation)
	}
}

// processApplyOperation processa operação de aplicação
func (p *Pool) processApplyOperation(item entities.QueueItem) (*entities.BoostOperation, error) {
	res, err := p.processor.ProcessApply(item.Context, item.BoosterID)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return &entities.BoostOperation{
		ID:        item.ID,
		BoosterID: item.BoosterID,
		Type:      item.Operation,
		AppliedAt: time.Now(),
		ErrorMsg:  p.getErrorMessage(res.Error),
	}, nil
}

func (p *Pool) getErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// processRevertOperation processa operação de reversão
func (p *Pool) processRevertOperation(item entities.QueueItem) (*entities.BoostOperation, error) {
	res, err := p.processor.ProcessRevert(item.Context, item.BoosterID)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return &entities.BoostOperation{
		ID:         item.ID,
		BoosterID:  item.BoosterID,
		Type:       item.Operation,
		RevertedAt: time.Now(),
		ErrorMsg:   p.getErrorMessage(res.Error),
	}, nil
}

// Stop para todos os workers
func (p *Pool) Stop() {
	p.stopOnce.Do(func() {
		p.queueManager.Stop()
		p.wg.Wait()
	})
}

func (p *Pool) GetActiveWorkerCount() int {
	return p.workerCount
}

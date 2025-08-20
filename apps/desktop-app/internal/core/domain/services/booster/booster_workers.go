package booster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type Processor interface {
	ProcessApply(ctx context.Context, boosterID string) (*entities.BoostApplyResult, error)
	ProcessRevert(ctx context.Context, boosterID string) (*entities.BoostRevertResult, error)
}

type EventEmitter interface {
	EmitProcessing(boosterID, operationID string, operation entities.BoosterOperationType)
	EmitSuccess(boosterID, operationID string, operation entities.BoosterOperationType, message string)
	EmitError(boosterID, operationID string, operation entities.BoosterOperationType, err error)
	EmitFailed(boosterID, operationID string, operation entities.BoosterOperationType, message string)
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
	p.queueManager.Remove(item.BoosterID)
	p.eventEmitter.EmitProcessing(item.BoosterID, item.OperationID, item.Operation)

	var op *entities.BoostOperation
	var err error

	switch item.Operation {
	case entities.BoosterOperationType(entities.ApplyOperationType):
		res, e := p.processor.ProcessApply(item.Context, item.BoosterID)
		err = e
		if res != nil {
			op = &entities.BoostOperation{
				ID:         item.ID,
				BoosterID:  item.BoosterID,
				Type:       item.Operation,
				AppliedAt: time.Now(),
				ErrorMsg:   res.Error.Error(),
			}
		}
	case entities.BoosterOperationType(entities.RevertOperationType):
		res, e := p.processor.ProcessRevert(item.Context, item.BoosterID)
		err = e
		if res != nil {
			op = &entities.BoostOperation{
				ID:         item.ID,
				BoosterID:  item.BoosterID,
				Type:       item.Operation,
				RevertedAt: time.Now(),
				ErrorMsg:   res.Error.Error(),
			}
		}
	default:
		err = fmt.Errorf("invalid operation")
	}

	// Emitir eventos baseados no status
	if err != nil {
		p.eventEmitter.EmitError(item.BoosterID, item.OperationID, item.Operation, err)
	} else {
		p.eventEmitter.EmitSuccess(item.BoosterID, item.OperationID, item.Operation, op.ErrorMsg)
	} 

	// Registrar no histórico
	p.historyRecorder.RecordOperation(item, op, err)
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

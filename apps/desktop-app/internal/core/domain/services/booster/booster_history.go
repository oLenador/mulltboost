// Package history gerencia o histórico de operações de boosters
package booster

import (
	"context"
	"time"

	repos "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/repositories"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

// Recorder gerencia a gravação do histórico de operações
type Recorder struct {
	operationsRepo *repos.BoostOperationsRepo
}

// NewRecorder cria um novo recorder de histórico
func NewRecorder(operationsRepo *repos.BoostOperationsRepo) *Recorder {
	return &Recorder{
		operationsRepo: operationsRepo,
	}
}

// RecordOperation registra uma operação no histórico
func (r *Recorder) RecordOperation(item entities.QueueItem, result *entities.BoostOperation, err error) error {
	operation := &entities.BoostOperation{
		ID:        item.OperationID,
		BoosterID: item.BoosterID,
		Type:      item.Operation,
		AppliedAt: item.SubmittedAt,
	}

	// save and ignore error
	if saveErr := r.operationsRepo.Save(context.Background(), operation); saveErr != nil {
		return err
	}
	return nil
}

func (r *Recorder) GetOperationsHistory(ctx context.Context, boosterID string) (*[]entities.BoostOperation, error) {
	return r.operationsRepo.GetByBoosterID(ctx, boosterID)
}

func (r *Recorder) GetAllOperations(ctx context.Context) (*[]entities.BoostOperation, error) {
	return r.operationsRepo.GetAll(ctx)
}

func (r *Recorder) GetOperationsByStatus(ctx context.Context, status entities.BoosterExecutionStatus) (*[]entities.BoostOperation, error) {
	return r.operationsRepo.GetByStatus(ctx, status)
}

func (r *Recorder) GetOperationsByType(ctx context.Context, operationType entities.BoosterOperationType) (*[]entities.BoostOperation, error) {
	return r.operationsRepo.GetByType(ctx, operationType)
}

func (r *Recorder) GetOperationsByTimeRange(ctx context.Context, startTime, endTime time.Time) (*[]entities.BoostOperation, error) {
	return r.operationsRepo.GetByTimeRange(ctx, startTime, endTime)
}

func (r *Recorder) GetRecentOperations(ctx context.Context, limit int) (*[]entities.BoostOperation, error) {
	return r.operationsRepo.GetRecent(ctx, limit)
}

func (r *Recorder) DeleteOldOperations(ctx context.Context, olderThan time.Time) error {
	return r.operationsRepo.DeleteOlderThan(ctx, olderThan)
}

func (r *Recorder) GetOperationStats(ctx context.Context) (*OperationStats, error) {
	allOps, err := r.GetAllOperations(ctx)
	if err != nil {
		return nil, err
	}

	if allOps == nil || len(*allOps) == 0 {
		return &OperationStats{}, nil
	}

	ops := *allOps
	stats := &OperationStats{
		Total: len(ops),
	}

	var totalDuration time.Duration
	var completedOps int

	for _, op := range ops {
		updateTypeStats(stats, op)
	}

	calculateRates(stats)
	calculateAverageDuration(stats, totalDuration, completedOps)

	return stats, nil
}

func updateTypeStats(stats *OperationStats, op entities.BoostOperation) {
	switch op.Type {
	case entities.BoosterOperationType(entities.ExecutionApplied):
		stats.ApplyOperations++
	case entities.BoosterOperationType(entities.RevertOperationType):
		stats.RevertOperations++
	}
}

func calculateRates(stats *OperationStats) {
	if stats.Total == 0 {
		return
	}
	stats.SuccessRate = float64(stats.Successful) / float64(stats.Total) * 100
	stats.FailureRate = float64(stats.Failed) / float64(stats.Total) * 100
}

func calculateAverageDuration(stats *OperationStats, totalDuration time.Duration, completedOps int) {
	if completedOps == 0 {
		return
	}
	stats.AverageDuration = totalDuration / time.Duration(completedOps)
}

type OperationStats struct {
	Total            int           `json:"total"`
	Successful       int           `json:"successful"`
	Failed           int           `json:"failed"`
	Processing       int           `json:"processing"`
	Pending          int           `json:"pending"`
	ApplyOperations  int           `json:"applyOperations"`
	RevertOperations int           `json:"revertOperations"`
	SuccessRate      float64       `json:"successRate"`
	FailureRate      float64       `json:"failureRate"`
	AverageDuration  time.Duration `json:"averageDuration"`
}
package booster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	repos "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/repositories"
)

type BoosterProcessor struct {
	rollbackRepo *repos.RollbackRepo
	boosters     map[string]inbound.BoosterUseCase
	boostersMu   sync.RWMutex
}

func NewBoosterProcessor(rollbackRepo *repos.RollbackRepo) *BoosterProcessor {
	return &BoosterProcessor{
		rollbackRepo: rollbackRepo,
		boosters:     make(map[string]inbound.BoosterUseCase),
	}
}

// RegisterBooster registra um novo booster
func (p *BoosterProcessor) RegisterBooster(booster inbound.BoosterUseCase) error {
	p.boostersMu.Lock()
	defer p.boostersMu.Unlock()
	
	info := booster.GetEntity()
	p.boosters[info.ID] = booster
	return nil
}

// GetBooster retorna um booster pelo ID
func (p *BoosterProcessor) GetBooster(id string) (inbound.BoosterUseCase, bool) {
	p.boostersMu.RLock()
	defer p.boostersMu.RUnlock()
	
	booster, exists := p.boosters[id]
	return booster, exists
}

// GetAllBoosters retorna todos os boosters registrados
func (p *BoosterProcessor) GetAllBoosters() map[string]inbound.BoosterUseCase {
	p.boostersMu.RLock()
	defer p.boostersMu.RUnlock()
	
	return p.boosters
}

func (p *BoosterProcessor) GetAllBoostersEntities() []entities.Booster {
	p.boostersMu.RLock()
	defer p.boostersMu.RUnlock()
	result := make([]entities.Booster, 0, len(p.boosters))

	for _, booster := range p.boosters {
		result = append(result, booster.GetEntity())
	}
	return result
}

// ProcessApply processa a aplicação de um booster
func (p *BoosterProcessor) ProcessApply(ctx context.Context, boosterID string) (*entities.BoostApplyResult, error) {
	booster, exists := p.GetBooster(boosterID)
	if !exists {
		return nil, fmt.Errorf("booster with ID %s not found", boosterID)
	}

	// Verifica se pode ser aplicado
	if !booster.CanApply(ctx) {
		return &entities.BoostApplyResult{
			Success: false,
			Message: "Booster cannot be applied at this time",
		}, nil
	}

	// Valida o booster
	if err := booster.Validate(ctx); err != nil {
		return &entities.BoostApplyResult{
			Success: false,
			Message: "Validation failed: " + err.Error(),
			Error:   err,
		}, nil
	}

	result, err := booster.Execute(ctx)
	if err != nil {
		return result, err
	}

	if err := p.saveRollbackState(ctx, boosterID, result, booster); err != nil {
		return result, fmt.Errorf("failed to save booster state: %w", err)
	}

	return result, nil
}

// ProcessRevert processa a reversão de um booster
func (p *BoosterProcessor) ProcessRevert(ctx context.Context, boosterID string) (*entities.BoostRevertResult, error) {
	booster, exists := p.GetBooster(boosterID)
	if !exists {
		return nil, fmt.Errorf("booster with ID %s not found", boosterID)
	}

	// Verifica se pode ser revertido
	if !booster.CanRevert(ctx) {
		return &entities.BoostRevertResult{
			Error: fmt.Errorf("booster cannot be reverted at this time"),
			Success: false,
			Message: "Booster cannot be reverted at this time",
		}, nil
	}

	// Executa a reversão
	result, err := booster.Revert(ctx)
	if err != nil {
		return nil, err
	}

	if err := p.updateRollbackState(ctx, boosterID, result); err != nil {
		// Log o erro mas não falha a operação
		fmt.Printf("Warning: failed to update rollback state for %s: %v\n", boosterID, err)
	}

	return result, nil
}

// saveRollbackState salva o estado de rollback após aplicação
func (p *BoosterProcessor) saveRollbackState(ctx context.Context, boosterID string, result *entities.BoostApplyResult, booster inbound.BoosterUseCase) error {
	state := &entities.BoosterRollbackState{
		ID:         boosterID,
		Applied:    result.Success,
		BackupData: result.BackupData,
		Version:    booster.GetEntity().Version,
	}

	if result.Success {
		now := time.Now()
		state.AppliedAt = &now
		state.Status = entities.ExecutionApplied
	} else {
		state.Status = entities.ExecutionFailed
		state.ErrorMsg = result.Message
	}

	return p.rollbackRepo.Save(ctx, state)
}

// updateRollbackState atualiza o estado de rollback após reversão
func (p *BoosterProcessor) updateRollbackState(ctx context.Context, boosterID string, result *entities.BoostRevertResult) error {
	rollbackEntity, err := p.rollbackRepo.GetByID(ctx, boosterID)
	if err != nil {
		return err
	}

	if rollbackEntity == nil {
		return fmt.Errorf("rollback state not found for booster %s", boosterID)
	}

	if result.Success {
		now := time.Now()
		rollbackEntity.RevertedAt = &now
		rollbackEntity.Status = entities.ExecutionReverted
		rollbackEntity.Applied = false
	} else {
		rollbackEntity.Status = entities.ExecutionFailed
		rollbackEntity.ErrorMsg = result.Message
	}

	return p.rollbackRepo.Save(ctx, rollbackEntity)
}

// GetRollbackState retorna o estado de rollback de um booster
func (p *BoosterProcessor) GetRollbackState(ctx context.Context, boosterID string) (*entities.BoosterRollbackState, error) {
	return p.rollbackRepo.GetByID(ctx, boosterID)
}

// ValidateBooster validate if a booster operationcan be applied
func (p *BoosterProcessor) ValidateBoosterOperation(ctx context.Context, boosterID string, operation entities.BoosterOperationType) error {
	booster, exists := p.GetBooster(boosterID)
	if !exists {
		return fmt.Errorf("booster with ID %s not found", boosterID)
	}

	switch operation {
	case entities.BoosterOperationType(entities.ApplyOperationType):
		if !booster.CanApply(ctx) {
			return fmt.Errorf("booster cannot be applied at this time")
		}
		return booster.Validate(ctx)
	case entities.BoosterOperationType(entities.RevertOperationType):
		if !booster.CanRevert(ctx) {
			return fmt.Errorf("booster cannot be reverted at this time")
		}
		return nil
	default:
		return fmt.Errorf("invalid operation type: %s", operation)
	}
}

// GetBoosterCount retorna o número de boosters registrados
func (p *BoosterProcessor) GetBoosterCount() int {
	p.boostersMu.RLock()
	defer p.boostersMu.RUnlock()
	return len(p.boosters)
}
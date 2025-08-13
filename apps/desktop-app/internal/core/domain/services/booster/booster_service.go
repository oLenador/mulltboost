package booster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/ports/outbound"
)

type Service struct {
	stateRepo  outbound.BoosterStateRepository
	boosters   map[string]inbound.BoosterUseCase
	boostersMu sync.RWMutex
}

func NewService(stateRepo outbound.BoosterStateRepository) *Service {
	return &Service{
		stateRepo: stateRepo,
		boosters:  make(map[string]inbound.BoosterUseCase),
	}
}

func (s *Service) RegisterBooster(booster inbound.BoosterUseCase) error {
	s.boostersMu.Lock()
	defer s.boostersMu.Unlock()

	info := booster.GetEntity()
	s.boosters[info.ID] = booster
	return nil
}

func (s *Service) GetAvailableBoosters(lang i18n.Language) []dto.BoosterDto {
	s.boostersMu.RLock()
	defer s.boostersMu.RUnlock()

	boosters := make([]dto.BoosterDto, 0, len(s.boosters))
	for _, booster := range s.boosters {
		boosters = append(boosters, booster.GetEntityDto(lang))
	}
	return boosters
}

func (s *Service) GetBoosterState(id string) (*entities.BoosterState, error) {
	return s.stateRepo.GetByID(context.Background(), id)
}

func (s *Service) ApplyBooster(ctx context.Context, id string) (*entities.BoosterResult, error) {
	s.boostersMu.RLock()
	booster, exists := s.boosters[id]
	s.boostersMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("booster with ID %s not found", id)
	}

	if !booster.CanApply(ctx) {
		return &entities.BoosterResult{
			Success: false,
			Message: "Booster cannot be applied at this time",
		}, nil
	}

	if err := booster.Validate(ctx); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: "Validation failed: " + err.Error(),
			Error:   err,
		}, nil
	}

	result, err := booster.Execute(ctx)
	if err != nil {
		return result, err
	}

	state := &entities.BoosterState{
		ID:         id,
		Applied:    result.Success,
		Status:     entities.StatusApplied,
		BackupData: result.BackupData,
		Version:    booster.GetEntity().Version,
	}

	if result.Success {
		now := time.Now()
		state.AppliedAt = &now
	} else {
		state.Status = entities.StatusFailed
		state.ErrorMsg = result.Message
	}

	if err := s.stateRepo.Save(ctx, state); err != nil {
		return result, fmt.Errorf("failed to save booster state: %w", err)
	}

	return result, nil
}

func (s *Service) RevertBooster(ctx context.Context, id string) (*entities.BoosterResult, error) {
	s.boostersMu.RLock()
	booster, exists := s.boosters[id]
	s.boostersMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("booster with ID %s not found", id)
	}

	if !booster.CanRevert(ctx) {
		return &entities.BoosterResult{
			Success: false,
			Message: "Booster cannot be reverted at this time",
		}, nil
	}

	result, err := booster.Revert(ctx)
	if err != nil {
		return result, err
	}

	state, _ := s.stateRepo.GetByID(ctx, id)
	if state != nil {
		if result.Success {
			now := time.Now()
			state.RevertedAt = &now
			state.Status = entities.StatusReverted
			state.Applied = false
		} else {
			state.Status = entities.StatusFailed
			state.ErrorMsg = result.Message
		}
		s.stateRepo.Save(ctx, state)
	}

	return result, nil
}

func (s *Service) ApplyBoosterBatch(ctx context.Context, ids []string) (*entities.BatchResult, error) {
	result := &entities.BatchResult{
		TotalCount: len(ids),
		Results:    make(map[string]entities.BoosterResult),
	}

	for _, id := range ids {
		boosterResult, err := s.ApplyBooster(ctx, id)
		if err != nil {
			result.Results[id] = entities.BoosterResult{
				Success: false,
				Message: err.Error(),
				Error:   err,
			}
			result.FailedCount++
		} else {
			result.Results[id] = *boosterResult
			if boosterResult.Success {
				result.SuccessCount++
			} else {
				result.FailedCount++
			}
		}
	}

	if result.SuccessCount == result.TotalCount {
		result.OverallStatus = "success"
	} else if result.FailedCount == result.TotalCount {
		result.OverallStatus = "failed"
	} else {
		result.OverallStatus = "partial"
	}

	return result, nil
}

func (s *Service) GetBoostersByCategory(category entities.BoosterCategory) []entities.Booster {
	s.boostersMu.RLock()
	defer s.boostersMu.RUnlock()

	var boosters []entities.Booster
	for _, booster := range s.boosters {
		info := booster.GetEntity()
		if info.Category == category {
			boosters = append(boosters, info)
		}
	}
	return boosters
}

package optimization

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/oLenador/mulltbost/internal/core/domain/entities"
    "github.com/oLenador/mulltbost/internal/core/ports/inbound"
    "github.com/oLenador/mulltbost/internal/core/ports/outbound"
)

type Service struct {
    stateRepo   outbound.OptimizationStateRepository
    plugins     map[string]inbound.OptimizationUseCase
    pluginsMux  sync.RWMutex
}

func NewService(stateRepo outbound.OptimizationStateRepository) *Service {
    return &Service{
        stateRepo: stateRepo,
        plugins:   make(map[string]inbound.OptimizationUseCase),
    }
}

func (s *Service) RegisterPlugin(plugin inbound.OptimizationUseCase) error {
    s.pluginsMux.Lock()
    defer s.pluginsMux.Unlock()
    
    info := plugin.GetInfo()
    s.plugins[info.ID] = plugin
    return nil
}

func (s *Service) GetAvailableOptimizations() []entities.Optimization {
    s.pluginsMux.RLock()
    defer s.pluginsMux.RUnlock()
    
    optimizations := make([]entities.Optimization, 0, len(s.plugins))
    for _, plugin := range s.plugins {
        optimizations = append(optimizations, plugin.GetInfo())
    }
    return optimizations
}

func (s *Service) GetOptimizationState(id string) (*entities.OptimizationState, error) {
    return s.stateRepo.GetByID(context.Background(), id)
}

func (s *Service) ApplyOptimization(ctx context.Context, id string) (*entities.OptimizationResult, error) {
    s.pluginsMux.RLock()
    plugin, exists := s.plugins[id]
    s.pluginsMux.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("optimization with ID %s not found", id)
    }

    // Verifica se pode aplicar
    if !plugin.CanApply(ctx) {
        return &entities.OptimizationResult{
            Success: false,
            Message: "Optimization cannot be applied at this time",
        }, nil
    }

    // Valida antes de aplicar
    if err := plugin.Validate(ctx); err != nil {
        return &entities.OptimizationResult{
            Success: false,
            Message: "Validation failed: " + err.Error(),
            Error:   err,
        }, nil
    }

    // Aplica a otimização
    result, err := plugin.Execute(ctx)
    if err != nil {
        return result, err
    }

    // Salva o estado
    state := &entities.OptimizationState{
        ID:         id,
        Applied:    result.Success,
        Status:     entities.StatusApplied,
        BackupData: result.BackupData,
        Version:    plugin.GetInfo().Version,
    }
    
    if result.Success {
        now := time.Now()
        state.AppliedAt = &now
    } else {
        state.Status = entities.StatusFailed
        state.ErrorMsg = result.Message
    }

    if err := s.stateRepo.Save(ctx, state); err != nil {
        return result, fmt.Errorf("failed to save optimization state: %w", err)
    }

    return result, nil
}

func (s *Service) RevertOptimization(ctx context.Context, id string) (*entities.OptimizationResult, error) {
    s.pluginsMux.RLock()
    plugin, exists := s.plugins[id]
    s.pluginsMux.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("optimization with ID %s not found", id)
    }

    if !plugin.CanRevert(ctx) {
        return &entities.OptimizationResult{
            Success: false,
            Message: "Optimization cannot be reverted at this time",
        }, nil
    }

    result, err := plugin.Revert(ctx)
    if err != nil {
        return result, err
    }

    // Atualiza o estado
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

func (s *Service) ApplyOptimizationBatch(ctx context.Context, ids []string) (*entities.BatchResult, error) {
    result := &entities.BatchResult{
        TotalCount: len(ids),
        Results:    make(map[string]entities.OptimizationResult),
    }

    for _, id := range ids {
        optResult, err := s.ApplyOptimization(ctx, id)
        if err != nil {
            result.Results[id] = entities.OptimizationResult{
                Success: false,
                Message: err.Error(),
                Error:   err,
            }
            result.FailedCount++
        } else {
            result.Results[id] = *optResult
            if optResult.Success {
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

func (s *Service) GetOptimizationsByCategory(category entities.OptimizationCategory) []entities.Optimization {
    s.pluginsMux.RLock()
    defer s.pluginsMux.RUnlock()
    
    var optimizations []entities.Optimization
    for _, plugin := range s.plugins {
        info := plugin.GetInfo()
        if info.Category == category {
            optimizations = append(optimizations, info)
        }
    }
    return optimizations
}
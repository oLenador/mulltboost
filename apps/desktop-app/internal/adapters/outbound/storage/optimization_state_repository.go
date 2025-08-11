package storage

import (
    "context"
    "sync"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type OptimizationStateRepository struct {
    states map[string]*entities.OptimizationState
    mutex  sync.RWMutex
}

func NewOptimizationStateRepository() *OptimizationStateRepository {
    return &OptimizationStateRepository{
        states: make(map[string]*entities.OptimizationState),
    }
}

func (r *OptimizationStateRepository) Save(ctx context.Context, state *entities.OptimizationState) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    r.states[state.ID] = state
    return nil
}

func (r *OptimizationStateRepository) GetByID(ctx context.Context, id string) (*entities.OptimizationState, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    state, exists := r.states[id]
    if !exists {
        return nil, nil
    }
    return state, nil
}

func (r *OptimizationStateRepository) GetAll(ctx context.Context) ([]*entities.OptimizationState, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    states := make([]*entities.OptimizationState, 0, len(r.states))
    for _, state := range r.states {
        states = append(states, state)
    }
    return states, nil
}

func (r *OptimizationStateRepository) Delete(ctx context.Context, id string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    delete(r.states, id)
    return nil
}
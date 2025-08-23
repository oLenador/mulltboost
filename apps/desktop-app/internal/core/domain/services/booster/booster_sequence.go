package booster

import (
	"sync"
	"time"
)

type SequenceManager struct {
	sequences map[string]*BoosterSequence
	mutex     sync.RWMutex
}

type BoosterSequence struct {
	CurrentSequence int64
	LastOperationID string
	LastEventTime   time.Time
	IsCompleted     bool
}

func NewSequenceManager() *SequenceManager {
	return &SequenceManager{
		sequences: make(map[string]*BoosterSequence),
	}
}

func (sm *SequenceManager) GetNextSequence(boosterID string, operationID string) int64 {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sequence, exists := sm.sequences[boosterID]
	if !exists || sequence.IsCompleted {
		sm.sequences[boosterID] = &BoosterSequence{
			CurrentSequence: 1,
			LastOperationID: operationID,
			LastEventTime:   time.Now(),
			IsCompleted:     false,
		}
		return 1
	}

	sequence.CurrentSequence++
	sequence.LastOperationID = operationID
	sequence.LastEventTime = time.Now()

	return sequence.CurrentSequence
}

func (sm *SequenceManager) MarkCompleted(boosterID string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if sequence, exists := sm.sequences[boosterID]; exists {
		sequence.IsCompleted = true
	}
}

func (sm *SequenceManager) GetCurrentSequence(boosterID string) int64 {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if sequence, exists := sm.sequences[boosterID]; exists {
		return sequence.CurrentSequence
	}
	return 0
}

func (sm *SequenceManager) IsIdempotent(boosterID string, operationID string) bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if sequence, exists := sm.sequences[boosterID]; exists {
		return sequence.LastOperationID == operationID && !sequence.IsCompleted
	}
	return false
}

func (sm *SequenceManager) CleanupOldSequences(maxAge time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	now := time.Now()
	for boosterID, sequence := range sm.sequences {
		if sequence.IsCompleted && now.Sub(sequence.LastEventTime) > maxAge {
			delete(sm.sequences, boosterID)
		}
	}
}

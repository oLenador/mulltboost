package booster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type Manager struct {
	items    []entities.QueueItem
	itemsMap map[string]*entities.QueueItem
	mu       sync.RWMutex
	workCh   chan entities.QueueItem
	stopCh   chan struct{}

	totalProcessed int
	inProgress     int
}

func NewManager(bufferSize int) *Manager {
	return &Manager{
		items:    make([]entities.QueueItem, 0),
		itemsMap: make(map[string]*entities.QueueItem),
		workCh:   make(chan entities.QueueItem, bufferSize),
		stopCh:   make(chan struct{}),
	}
}

// Add adiciona um item à queue, removendo duplicatas e conflitos
func (m *Manager) Add(boosterID string, operation entities.BoosterOperationType) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Se já existe uma operação para este booster
	if existingItem, exists := m.itemsMap[boosterID]; exists {
		// Se é a mesma operação, retorna o ID existente
		if existingItem.Operation == operation {
			return existingItem.OperationID, nil
		}

		// Se é operação contrária, cancela a existente
		existingItem.Cancel()
		m.removeUnsafe(boosterID)
	}

	// Cria novo item
	ctx, cancel := context.WithCancel(context.Background())
	operationID := uuid.New().String()

	item := entities.QueueItem{
		ID:          uuid.New().String(),
		BoosterID:   boosterID,
		Operation:   operation,
		OperationID: operationID,
		SubmittedAt: time.Now(),
		Context:     ctx,
		Cancel:      cancel,
	}

	// Adiciona à queue
	m.items = append(m.items, item)
	m.itemsMap[boosterID] = &m.items[len(m.items)-1]

	m.inProgress++

	// Envia para processamento
	select {
	case m.workCh <- item:
		return operationID, nil
	default:
		// Queue cheia, remove item e retorna erro
		m.removeUnsafe(boosterID)
		cancel()
		return "", ErrQueueFull
	}
}

// Remove remove um item da queue
func (m *Manager) Remove(boosterID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.removeUnsafe(boosterID)
}

// removeUnsafe remove um item da queue (sem lock)
func (m *Manager) removeUnsafe(boosterID string) {
	if item, exists := m.itemsMap[boosterID]; exists {
		// Remove do slice
		for i, qItem := range m.items {
			if qItem.BoosterID == boosterID {
				m.items = append(m.items[:i], m.items[i+1:]...)
				break
			}
		}

		// Remove do map
		delete(m.itemsMap, boosterID)

		if item.Cancel != nil {
			item.Cancel()
		}

		// Atualiza estatísticas
		m.inProgress--
		m.totalProcessed++

		// Se a fila está vazia, zera tudo
		if len(m.items) == 0 {
			m.inProgress = 0
			m.totalProcessed = 0
		}
	}
}
func (m *Manager) GetQueueStats() *entities.QueueState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]entities.QueueItem, len(m.items))
	copy(items, m.items)
	return &entities.QueueState{
		Items:          items,
		QueueSize:      m.totalProcessed + m.inProgress,
		TotalProcessed: m.totalProcessed,
		InProgress:     m.inProgress,
	}
}

// GetPosition retorna a posição na queue
func (m *Manager) GetPosition(boosterID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for i, item := range m.items {
		if item.BoosterID == boosterID {
			return i + 1
		}
	}
	return -1
}

// IsInQueue verifica se um booster está na queue
func (m *Manager) IsInQueue(boosterID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.itemsMap[boosterID]
	return exists
}

func (m *Manager) GetQueuedItem(boosterID string) *entities.QueueItem {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.itemsMap[boosterID]
	if !exists {
		return nil
	}

	// Retorna uma cópia para evitar race conditions
	itemCopy := *item
	return &itemCopy
}

// GetAllItems retorna uma cópia de todos os itens na queue
func (m *Manager) GetAllItems() []entities.QueueItem {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := make([]entities.QueueItem, len(m.items))
	copy(items, m.items)
	return items
}

// GetWorkChannel retorna o canal de trabalho para os workers
func (m *Manager) GetWorkChannel() <-chan entities.QueueItem {
	return m.workCh
}

// GetStopChannel retorna o canal de parada
func (m *Manager) GetStopChannel() <-chan struct{} {
	return m.stopCh
}

// Size retorna o tamanho atual da queue
func (m *Manager) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.items)
}

// Clear limpa toda a queue, cancelando todos os itens
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Cancela todos os contextos
	for _, item := range m.items {
		if item.Cancel != nil {
			item.Cancel()
		}
	}

	// Limpa as estruturas
	m.items = make([]entities.QueueItem, 0)
	m.itemsMap = make(map[string]*entities.QueueItem)
}

// Stop para o gerenciador e limpa recursos
func (m *Manager) Stop() {
	close(m.stopCh)
	m.Clear()
}

var (
	ErrQueueFull = fmt.Errorf("execution queue is full")
	ErrNotFound  = fmt.Errorf("item not found in queue")
)

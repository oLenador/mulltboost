package booster

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	boosterBase "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection"
	repos "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/repositories"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	processor           *BoosterProcessor
	queueManager        *Manager
	workerPool          *Pool
	historyRecorder     *Recorder
	eventEmitter        *BoosterEventEmitter
	boostActivationRepo *repos.BoostConfigRepository
}

type Config struct {
	WorkerCount     int
	QueueBufferSize int
}

func NewService(
	rollbackRepo *repos.RollbackRepo,
	operationsRepo *repos.BoostOperationsRepo,
	eventManager *application.EventManager,
	boostActivationRepo *repos.BoostConfigRepository,
) (*Service, error) {
	config := Config{
		WorkerCount:     3,
		QueueBufferSize: 100,
	}

	boosterProcessor := NewBoosterProcessor(rollbackRepo)
	queueManager := NewManager(config.QueueBufferSize)
	historyRecorder := NewRecorder(operationsRepo)
	eventEmitter := NewBoosterEventEmitter(eventManager)

	workerPool := NewPool(
		config.WorkerCount,
		boosterProcessor,
		eventEmitter,
		historyRecorder,
		queueManager,
	)

	err := initAllBoosts(boosterProcessor)
	if err != nil {
		return nil, fmt.Errorf("Erro on register the boosters", err)
	}

	err = boostActivationRepo.SyncWithAvailableBoosts(boosterProcessor.GetAllBoostersEntities())
	if err != nil {
		return nil, fmt.Errorf("Erro on sync the boosters", err)
	}

	service := &Service{
		processor:           boosterProcessor,
		queueManager:        queueManager,
		workerPool:          workerPool,
		historyRecorder:     historyRecorder,
		eventEmitter:        eventEmitter,
		boostActivationRepo: boostActivationRepo,
	}

	service.StartWorkers()

	return service, nil
}

func initAllBoosts(processor *BoosterProcessor) error {

	ps := boosterBase.GetPlatformServices()
	deps := inbound.NewExecutorDepServices(ps)

	loaders := map[string][]inbound.BoosterUseCase{
		"connection": connection.GetAllPlugins(deps),
	}

	for _, boostArray := range loaders {
		for _, booster := range boostArray {
			if err := processor.RegisterBooster(booster); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) StartWorkers() {
	s.workerPool.Start()
}

func (s *Service) StopWorkers() {
	s.workerPool.Stop()
}

func (s *Service) RegisterBooster(booster inbound.BoosterUseCase) error {
	return s.processor.RegisterBooster(booster)
}

func (s *Service) GetOperationsHistory(ctx context.Context, id string) (*[]entities.BoostOperation, error) {
	return s.historyRecorder.GetOperationsHistory(ctx, id)
}

func (s *Service) GetAllBoosters(ctx context.Context, lang i18n.Language) []entities.Booster {
	boosters := s.processor.GetAllBoosters()
	result := make([]entities.Booster, 0, len(boosters))

	for _, booster := range boosters {
		result = append(result, booster.GetEntity())
	}

	return result
}

// Método auxiliar para converter booster em DTO com estado de ativação
func (s *Service) convertBoosterToDto(ctx context.Context, booster inbound.BoosterUseCase, lang i18n.Language) (*dto.GetBoosterDto, error) {
	boosterDto := booster.GetEntityDto(lang)
	resActivationState, err := s.boostActivationRepo.GetBoostState(ctx, boosterDto.ID)
	if err != nil {
		return nil, err
	}
	
	return &dto.GetBoosterDto{
		ID:           boosterDto.ID,
		Name:         boosterDto.Name,
		Description:  boosterDto.Description,
		Category:     boosterDto.Category,
		Level:        boosterDto.Level,
		Platform:     boosterDto.Platform,
		Dependencies: boosterDto.Dependencies,
		Conflicts:    boosterDto.Conflicts,
		Reversible:   boosterDto.Reversible,
		RiskLevel:    boosterDto.RiskLevel,
		Version:      boosterDto.Version,
		IsApplied:    resActivationState.IsApplied,
		AppliedAt:    resActivationState.AppliedAt,
		RevertedAt:   resActivationState.RevertedAt,
		Tags:         boosterDto.Tags,
	}, nil
}


func (s *Service) GetAvailableBoosters(ctx context.Context, lang i18n.Language) []dto.GetBoosterDto {
	boosters := s.processor.GetAllBoosters()
	result := make([]dto.GetBoosterDto, 0, len(boosters))
	
	for _, booster := range boosters {
		boosterDto, err := s.convertBoosterToDto(ctx, booster, lang)
		if err != nil {
			continue
		}
		result = append(result, *boosterDto)
	}
	return result
}

func (s *Service) GetBoostersByCategory(ctx context.Context, category entities.BoosterCategory, lang i18n.Language) []dto.GetBoosterDto {
	boosters := s.processor.GetAllBoosters()
	logger.NewCustomLogger("GetBoosters").DebugFields(
		"Listando boosters",
		logger.Fields{
			"boosters": boosters,
		},
	)
	
	result := make([]dto.GetBoosterDto, 0)
	for _, booster := range boosters {
		boosterDto := booster.GetEntityDto(lang)
		if boosterDto.Category == category {
			fullBoosterDto, err := s.convertBoosterToDto(ctx, booster, lang)
			if err != nil {
				continue
			}
			result = append(result, *fullBoosterDto)
		}
	}
	return result
}

func (s *Service) GetBoosterQueueStatus(ctx context.Context, id string, lang i18n.Language) (*entities.QueueItem, error) {
	_, exists := s.processor.GetBooster(id)
	if !exists {
		return nil, fmt.Errorf("booster with ID %s not found", id)
	}
	return s.queueManager.GetQueuedItem(id), nil
}

func (s *Service) GetExecutionQueueState(ctx context.Context) *entities.QueueState {
	queueItems := s.queueManager.GetQueueStats()
	return queueItems
}

func (s *Service) InitBoosterApply(ctx context.Context, id string) (entities.InitResult, error) {
	if err := s.processor.ValidateBoosterOperation(ctx, id, entities.ApplyOperationType); err != nil {
		return entities.InitResult{
			OperationID: "",
			SubmittedAt: time.Now(),
			Success:     false,
			Status:      entities.OperationFailed,
			Message:     "validation failed",
			Error:       err,
		}, err
	}

	operationID, err := s.queueManager.Add(id, entities.ApplyOperationType)
	if err != nil {
		return entities.InitResult{
			OperationID: "",
			SubmittedAt: time.Now(),
			Success:     false,
			Status:      entities.OperationFailed,
			Message:     "failed to enqueue operation",
			Error:       err,
		}, err
	}

	s.eventEmitter.EmitQueued(id, operationID, entities.ApplyOperationType, s.queueManager.Size())

	return entities.InitResult{
		OperationID: operationID,
		SubmittedAt: time.Now(),
		Success:     true,
		Status:      entities.OperationPending,
		Message:     "operation queued successfully",
	}, nil
}

func (s *Service) InitBoosterApplyBatch(ctx context.Context, ids []string) (entities.InitResult, error) {
	batchID := uuid.New().String()
	successCount := 0
	validationErrors := make(map[string]error)

	for _, id := range ids {
		if err := s.processor.ValidateBoosterOperation(ctx, id, entities.ApplyOperationType); err != nil {
			validationErrors[id] = err
			continue
		}

		if _, err := s.queueManager.Add(id, entities.ApplyOperationType); err == nil {
			successCount++
		} else {
			validationErrors[id] = err
		}
	}

	s.eventEmitter.EmitBatchQueued(
		batchID,
		entities.ApplyOperationType,
		len(ids),
		successCount,
		validationErrors,
		s.queueManager.Size(),
	)

	status := entities.OperationPending
	success := true
	message := fmt.Sprintf("batch queued with %d/%d successes", successCount, len(ids))

	if successCount == 0 {
		status = entities.OperationFailed
		success = false
		message = "all operations failed"
	}

	return entities.InitResult{
		OperationID: batchID,
		SubmittedAt: time.Now(),
		Success:     success,
		Status:      status,
		Message:     message,
	}, nil
}

func (s *Service) InitRevertBooster(ctx context.Context, id string) (entities.InitResult, error) {
	if err := s.processor.ValidateBoosterOperation(ctx, id, entities.RevertOperationType); err != nil {
		return entities.InitResult{
			SubmittedAt: time.Now(),
			Success:     false,
			Status:      entities.OperationFailed,
			Message:     "validation failed",
			Error:       err,
		}, err
	}

	operationID, err := s.queueManager.Add(id, entities.RevertOperationType)
	if err != nil {
		return entities.InitResult{
			SubmittedAt: time.Now(),
			Success:     false,
			Status:      entities.OperationFailed,
			Message:     "failed to enqueue operation",
			Error:       err,
		}, err
	}

	s.eventEmitter.EmitQueued(id, operationID, entities.RevertOperationType, s.queueManager.Size())

	return entities.InitResult{
		OperationID: operationID,
		SubmittedAt: time.Now(),
		Success:     true,
		Status:      entities.OperationPending,
		Message:     "revert operation queued successfully",
	}, nil
}

func (s *Service) InitRevertBoosterBatch(ctx context.Context, ids []string) (entities.InitResult, error) {
	batchID := uuid.New().String()
	successCount := 0
	validationErrors := make(map[string]error)

	for _, id := range ids {
		if err := s.processor.ValidateBoosterOperation(ctx, id, entities.RevertOperationType); err != nil {
			validationErrors[id] = err
			continue
		}

		if _, err := s.queueManager.Add(id, entities.RevertOperationType); err == nil {
			successCount++
		} else {
			validationErrors[id] = err
		}
	}

	s.eventEmitter.EmitBatchQueued(
		batchID,
		entities.RevertOperationType,
		len(ids),
		successCount,
		validationErrors,
		s.queueManager.Size(),
	)

	status := entities.OperationPending
	success := true
	message := fmt.Sprintf("batch queued with %d/%d successes", successCount, len(ids))

	if successCount == 0 {
		status = entities.OperationFailed
		success = false
		message = "all operations failed"
	}

	return entities.InitResult{
		OperationID: batchID,
		SubmittedAt: time.Now(),
		Success:     success,
		Status:      status,
		Message:     message,
	}, nil
}

func (s *Service) GetBoosterRollbackState(id string) (*entities.BoosterRollbackState, error) {
	return s.processor.GetRollbackState(context.Background(), id)
}

// CancelOperation cancela uma operação na queue
func (s *Service) CancelOperation(ctx context.Context, boosterID string) error {
	if !s.queueManager.IsInQueue(boosterID) {
		return fmt.Errorf("booster %s is not in queue", boosterID)
	}

	s.queueManager.Remove(boosterID)

	// Emite evento de cancelamento
	s.eventEmitter.EmitCancelled(boosterID, s.queueManager.Size())

	return nil
}

// GetQueueStats retorna estatísticas da queue
func (s *Service) GetQueueStats() *QueueStats {
	return &QueueStats{
		Size:          s.queueManager.Size(),
		ActiveWorkers: s.workerPool.GetActiveWorkerCount(),
		IsHealthy:     s.queueManager.Size() < 50, // Arbitrário
	}
}

func (s *Service) GetOperationStats(ctx context.Context) (*OperationStats, error) {
	return s.historyRecorder.GetOperationStats(ctx)
}

type QueueStats struct {
	Size          int  `json:"size"`
	ActiveWorkers int  `json:"activeWorkers"`
	IsHealthy     bool `json:"isHealthy"`
}

// HealthCheck verifica a saúde do serviço
func (s *Service) HealthCheck() *HealthStatus {
	queueSize := s.queueManager.Size()

	status := &HealthStatus{
		IsHealthy:          true,
		QueueSize:          queueSize,
		ActiveWorkers:      s.workerPool.GetActiveWorkerCount(),
		RegisteredBoosters: s.processor.GetBoosterCount(),
	}

	// Verifica se a queue não está muito cheia
	if queueSize > 80 {
		status.IsHealthy = false
		status.Issues = append(status.Issues, "Queue is nearly full")
	}

	return status
}

// HealthStatus contém informações sobre a saúde do serviço
type HealthStatus struct {
	IsHealthy          bool     `json:"isHealthy"`
	QueueSize          int      `json:"queueSize"`
	ActiveWorkers      int      `json:"activeWorkers"`
	RegisteredBoosters int      `json:"registeredBoosters"`
	Issues             []string `json:"issues,omitempty"`
}

package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	storage "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/models"
	"gorm.io/gorm"
)

type BoostConfigRepository struct {
	db             *gorm.DB
	cache          map[string]*storage.BoostActivationState
	lastSync       time.Time
	currentVersion string
}

const OBSOLETE_OLDER_THAN_DAYS int = 7

func NewBoostConfigRepository(db *gorm.DB) *BoostConfigRepository {
	repo := &BoostConfigRepository{
		db:             db,
		cache:          make(map[string]*storage.BoostActivationState),
		currentVersion: "1.0.0",
	}

	return repo
}

func (r *BoostConfigRepository) SyncWithAvailableBoosts(boosts []entities.Booster) error {
	// Criar mapa de boosts disponíveis
	availableBoosts := make(map[string]entities.Booster)
	for _, boost := range boosts {
		availableBoosts[boost.ID] = boost
	}

	// 1. Buscar estados existentes
	var existingStates []storage.BoostActivationState
	if err := r.db.Find(&existingStates).Error; err != nil {
		return fmt.Errorf("erro ao buscar estados existentes: %w", err)
	}
	existingMap := make(map[string]*storage.BoostActivationState)
	for i := range existingStates {
		// usar ponteiro do slice
		existingMap[existingStates[i].ID] = &existingStates[i]
	}

	// 2. Processar boosts disponíveis
	for boostKey, boost := range availableBoosts {
		if existing := existingMap[boostKey]; existing != nil {
			// Atualizar versão se necessário (adaptação: use a versão do boost)
			if existing.Version != boost.Version {
				existing.Version = boost.Version
				existing.UpdatedAt = time.Now()
				_ = r.db.Save(existing)
			}
			r.cache[boostKey] = existing
		} else {
			// Boost novo - criar estado
			newState := &storage.BoostActivationState{
				ID:        boostKey,
				IsApplied: false,
				Version:   boost.Version,
				Status:    entities.StatusInactive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := r.db.Create(newState).Error; err != nil {
				log.Printf("Erro ao criar estado para boost %s: %v", boostKey, err)
				continue
			}
			r.cache[boostKey] = newState
		}
	}

	// 3. Marcar boosts órfãos como obsoletos em vez de deletar
	for dbBoostKey, existing := range existingMap {
		if _, stillExists := availableBoosts[dbBoostKey]; !stillExists {
			// Boost removido do código - marcar como obsoleto e desativar
			if existing.IsApplied {
				log.Printf("Boost %s foi removido do código mas ainda está aplicado. Revertendo/Desativando...", dbBoostKey)
				existing.IsApplied = false
				existing.RevertedAt = timePtr(time.Now())
				existing.Status = entities.StatusObsolete
				existing.ErrorMessage = "Boost removido da versão atual"
				existing.UpdatedAt = time.Now()
				_ = r.db.Save(existing)
			} else {
				// Se já estava inativo, apenas marcar obsoleto e salvar
				existing.Status = entities.StatusObsolete
				existing.UpdatedAt = time.Now()
				_ = r.db.Save(existing)
			}
			// NÃO adicionar ao cache (fica invisível para a aplicação)
		}
	}

	r.lastSync = time.Now()
	return nil
}

func (r *BoostConfigRepository) GetAllActiveBoosts(ctx context.Context, boosts []entities.Booster) (map[string]*storage.BoostActivationState, error) {
	// Recarregar cache se necessário (a cada 1 hora)
	if time.Since(r.lastSync) > time.Hour {
		if err := r.SyncWithAvailableBoosts(boosts); err != nil {
			return nil, err
		}
	}

	result := make(map[string]*storage.BoostActivationState)
	for key, state := range r.cache {
		if state != nil && state.IsApplied {
			result[key] = state
		}
	}
	return result, nil
}

func (r *BoostConfigRepository) IsBoostActive(ctx context.Context, boostKey string) bool {
	if state, exists := r.cache[boostKey]; exists {
		return state.IsApplied
	}
	return false
}

func (r *BoostConfigRepository) ActivateBoost(ctx context.Context, boostKey string) error {
	state, exists := r.cache[boostKey]
	if !exists {
		return fmt.Errorf("boost %s não encontrado", boostKey)
	}
	if state.IsApplied {
		return fmt.Errorf("boost %s já está ativo", boostKey)
	}

	// Verificar se pode ser aplicado (rollback state, etc.)
	if canApply, err := r.canApplyBoost(ctx, boostKey); err != nil || !canApply {
		return fmt.Errorf("não é possível ativar boost %s: %v", boostKey, err)
	}

	now := time.Now()
	state.IsApplied = true
	state.AppliedAt = &now
	state.RevertedAt = nil
	state.Status = entities.StatusActive
	state.ErrorMessage = ""
	state.UpdatedAt = now

	if err := r.db.Save(state).Error; err != nil {
		return fmt.Errorf("erro ao ativar boost: %w", err)
	}

	return nil
}

func (r *BoostConfigRepository) DeactivateBoost(ctx context.Context, boostKey string) error {
	state, exists := r.cache[boostKey]
	if !exists {
		return fmt.Errorf("boost %s não encontrado", boostKey)
	}
	if !state.IsApplied {
		return fmt.Errorf("boost %s já está inativo", boostKey)
	}

	now := time.Now()
	state.IsApplied = false
	state.RevertedAt = &now
	state.Status = entities.StatusInactive
	state.UpdatedAt = now

	if err := r.db.Save(state).Error; err != nil {
		return fmt.Errorf("erro ao desativar boost: %w", err)
	}

	return nil
}

func (r *BoostConfigRepository) canApplyBoost(ctx context.Context, boostKey string) (bool, error) {
	// Verificar no BoosterRollbackState se já foi aplicado
	var rollbackState storage.BoosterRollbackState
	err := r.db.Where("id = ?", boostKey).First(&rollbackState).Error
	if err == gorm.ErrRecordNotFound {
		return true, nil // Nunca foi aplicado, pode aplicar
	}
	if err != nil {
		return false, err
	}
	// Se já foi aplicado e não foi revertido, não pode aplicar novamente
	if rollbackState.Applied && rollbackState.RevertedAt == nil {
		return false, fmt.Errorf("boost já foi aplicado anteriormente")
	}
	return true, nil
}

func (r *BoostConfigRepository) CleanupObsoleteStates(ctx context.Context) error {
	cutoff := time.Now().AddDate(0, 0, -OBSOLETE_OLDER_THAN_DAYS)
	result := r.db.Where("status = ? AND is_applied = ? AND updated_at < ?", entities.StatusObsolete, false, cutoff).Delete(&storage.BoostActivationState{})
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Removidos %d estados obsoletos de boosts", result.RowsAffected)
	return nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func (r *BoostConfigRepository) SetCurrentVersion(version string) {
	r.currentVersion = version
}

func (r *BoostConfigRepository) GetBoostState(ctx context.Context, boostKey string) (*storage.BoostActivationState, error) {
	if state, exists := r.cache[boostKey]; exists {
		return state, nil
	}
	// Se não está no cache, buscar no banco
	var state storage.BoostActivationState
	if err := r.db.Where("id = ?", boostKey).First(&state).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("boost %s não encontrado", boostKey)
		}
		return nil, err
	}
	// Adicionar ao cache
	r.cache[boostKey] = &state
	return &state, nil
}

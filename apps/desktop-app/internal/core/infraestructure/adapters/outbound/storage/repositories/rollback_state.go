package storage

import (
	"context"
	"errors"
	"fmt"

	mapper "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/mapper"
	storage "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/models"
	
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"gorm.io/gorm"
)

type RollbackRepo struct {
	db *gorm.DB
}

func NewRollbackRepo(db *gorm.DB) *RollbackRepo { return &RollbackRepo{db: db} }

func (r *RollbackRepo) Save(ctx context.Context, s *entities.BoosterRollbackState) error {
	if s == nil {
		return errors.New("nil rollback state")
	}
	model := mapper.MapRollbackFromDomain(s)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *RollbackRepo) GetByID(ctx context.Context, id string) (*entities.BoosterRollbackState, error) {
	var model storage.BoosterRollbackState
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return mapper.MapRollbackToDomain(&model), err
}

func (r *RollbackRepo) GetAll(ctx context.Context) ([]entities.BoosterRollbackState, error) {
	var models []storage.BoosterRollbackState
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.BoosterRollbackState, len(models))
	for i := range models {
		result[i] = *mapper.MapRollbackToDomain(&models[i])
	}
	return result, nil
}

func (r *RollbackRepo) GetByStatus(ctx context.Context, status entities.ExecutionStatus) ([]entities.BoosterRollbackState, error) {
	var models []storage.BoosterRollbackState
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at desc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.BoosterRollbackState, len(models))
	for i := range models {
		result[i] = *mapper.MapRollbackToDomain(&models[i])
	}
	return result, nil
}

func (r *RollbackRepo) UpdateStatus(ctx context.Context, id string, status entities.ExecutionStatus, errMsg string) error {
	updates := map[string]interface{}{
		"status":    storage.ExecutionStatus(status),
		"error_msg": errMsg,
	}

	if status == entities.StatusApplied {
		updates["applied"] = true
		updates["applied_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	}
	if status == entities.StatusReverted {
		updates["applied"] = false
		updates["reverted_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	}

	res := r.db.WithContext(ctx).Model(&storage.BoosterRollbackState{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("nenhuma linha afetada, id não encontrado: %s", id)
	}
	return nil
}

func (r *RollbackRepo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&storage.BoosterRollbackState{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("nenhuma linha afetada, id não encontrado: %s", id)
	}
	return nil
}
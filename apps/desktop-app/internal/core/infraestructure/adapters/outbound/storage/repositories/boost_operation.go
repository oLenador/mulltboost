package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	mapper "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/mapper"
	storage "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/models"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"gorm.io/gorm"
)

type BoostOperationsRepo struct {
	db *gorm.DB
}

func NewBoostOperationsRepo(db *gorm.DB) *BoostOperationsRepo { return &BoostOperationsRepo{db: db} }

func (r *BoostOperationsRepo) Save(ctx context.Context, entity *entities.BoostOperation) error {
	if entity == nil {
		return errors.New("nil applied boost")
	}
	model := mapper.MapAppliedFromDomain(entity)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *BoostOperationsRepo) GetByID(ctx context.Context, id string) (*entities.BoostOperation, error) {
	var a storage.BoostOperation
	err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return mapper.MapAppliedToDomain(&a), err
}

func (r *BoostOperationsRepo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&storage.BoostOperation{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("nenhuma linha afetada, id não encontrado: %s", id)
	}
	return nil
}

func (r *BoostOperationsRepo) GetByBoosterID(ctx context.Context, boosterID string) (*[]entities.BoostOperation, error) {
	var models []storage.BoostOperation
	if err := r.db.WithContext(ctx).Where(&storage.BoostOperation{BoosterID: boosterID}).Find(&models).Error; err != nil {
		return nil, err
	}
	domain := make([]entities.BoostOperation, 0, len(models))
	for i := range models {
		if m := mapper.MapAppliedToDomain(&models[i]); m != nil {
			domain = append(domain, *m)
		}
	}
	return &domain, nil
}

func (r *BoostOperationsRepo) GetAll(ctx context.Context) (*[]entities.BoostOperation, error) {
	var models []storage.BoostOperation
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	domain := make([]entities.BoostOperation, 0, len(models))
	for i := range models {
		if m := mapper.MapAppliedToDomain(&models[i]); m != nil {
			domain = append(domain, *m)
		}
	}
	return &domain, nil
}

func (r *BoostOperationsRepo) GetByStatus(ctx context.Context, status entities.BoosterExecutionStatus) (*[]entities.BoostOperation, error) {
	var models []storage.BoostOperation
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&models).Error; err != nil {
		return nil, err
	}
	domain := make([]entities.BoostOperation, 0, len(models))
	for i := range models {
		if m := mapper.MapAppliedToDomain(&models[i]); m != nil {
			domain = append(domain, *m)
		}
	}
	return &domain, nil
}

func (r *BoostOperationsRepo) GetByType(ctx context.Context, operationType entities.BoosterOperationType) (*[]entities.BoostOperation, error) {
	var models []storage.BoostOperation
	if err := r.db.WithContext(ctx).Where("type = ?", operationType).Find(&models).Error; err != nil {
		return nil, err
	}
	domain := make([]entities.BoostOperation, 0, len(models))
	for i := range models {
		if m := mapper.MapAppliedToDomain(&models[i]); m != nil {
			domain = append(domain, *m)
		}
	}
	return &domain, nil
}

func (r *BoostOperationsRepo) GetByTimeRange(ctx context.Context, startTime, endTime time.Time) (*[]entities.BoostOperation, error) {
	var models []storage.BoostOperation
	if err := r.db.WithContext(ctx).Where("applied_at BETWEEN ? AND ?", startTime, endTime).Find(&models).Error; err != nil {
		return nil, err
	}
	domain := make([]entities.BoostOperation, 0, len(models))
	for i := range models {
		if m := mapper.MapAppliedToDomain(&models[i]); m != nil {
			domain = append(domain, *m)
		}
	}
	return &domain, nil
}

func (r *BoostOperationsRepo) GetRecent(ctx context.Context, limit int) (*[]entities.BoostOperation, error) {
	if limit <= 0 {
		return &[]entities.BoostOperation{}, nil
	}
	var models []storage.BoostOperation
	if err := r.db.WithContext(ctx).Order("applied_at desc").Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}
	domain := make([]entities.BoostOperation, 0, len(models))
	for i := range models {
		if m := mapper.MapAppliedToDomain(&models[i]); m != nil {
			domain = append(domain, *m)
		}
	}
	return &domain, nil
}

func (r *BoostOperationsRepo) DeleteOlderThan(ctx context.Context, olderThan time.Time) error {
	res := r.db.WithContext(ctx).Where("applied_at < ?", olderThan).Delete(&storage.BoostOperation{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

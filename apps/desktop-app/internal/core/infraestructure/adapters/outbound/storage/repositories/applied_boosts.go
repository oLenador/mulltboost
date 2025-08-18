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

type AppliedRepo struct {
	db *gorm.DB
}

func NewAppliedRepo(db *gorm.DB) *AppliedRepo { return &AppliedRepo{db: db} }

func (r *AppliedRepo) Save(ctx context.Context, entity *entities.AppliedBoost) error {
	if entity == nil {
		return errors.New("nil applied boost")
	}
	
	model := mapper.MapAppliedFromDomain(entity)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *AppliedRepo) GetByID(ctx context.Context, id string) (*entities.AppliedBoost, error) {
	var a storage.AppliedBoost
	err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return mapper.MapAppliedToDomain(&a), err
}

func (r *AppliedRepo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&storage.AppliedBoost{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("nenhuma linha afetada, id não encontrado: %s", id)
	}
	return nil
}

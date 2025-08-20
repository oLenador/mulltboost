package storage

import (
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	storage "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/models"
)

func MapAppliedToDomain(a *storage.BoostOperation) *entities.BoostOperation {
	if a == nil {
		return nil
	}

	return &entities.BoostOperation{
		ID:         a.ID,
		BoosterID:  a.BoosterID,
		AppliedAt:  a.AppliedAt,
		RevertedAt: a.RevertedAt,
		Type:       a.Type,
		ErrorMsg:   a.ErrorMsg,
	}
}

func MapAppliedFromDomain(e *entities.BoostOperation) *storage.BoostOperation {
	if e == nil {
		return nil
	}

	return &storage.BoostOperation{
		ID:         e.ID,
		BoosterID:  e.BoosterID,
		AppliedAt:  e.AppliedAt,
		RevertedAt: e.RevertedAt,
		ErrorMsg:   e.ErrorMsg,
	}
}

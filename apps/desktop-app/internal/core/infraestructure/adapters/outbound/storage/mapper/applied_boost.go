package storage

import (
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	storage "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/models"
)


func MapAppliedToDomain(a *storage.AppliedBoost) *entities.AppliedBoost {
	if a == nil {
		return nil
	}



	return &entities.AppliedBoost{
		ID:         a.ID,
		UserID:     a.UserID,
		BoosterID:  a.BoosterID,
		Version:    a.Version,
		AppliedAt:  a.AppliedAt,
		RevertedAt: a.RevertedAt,
		Status:     entities.ExecutionStatus(a.Status),
		ErrorMsg:   a.ErrorMsg,
	}
}

func MapAppliedFromDomain(e *entities.AppliedBoost) *storage.AppliedBoost {
	if e == nil {
		return nil
	}


	return &storage.AppliedBoost{
		ID:         e.ID,
		UserID:     e.UserID,
		BoosterID:  e.BoosterID,
		Version:    e.Version,
		AppliedAt:  e.AppliedAt,
		RevertedAt: e.RevertedAt,
		Status:     storage.ExecutionStatus(e.Status),
		ErrorMsg:   e.ErrorMsg,
	}
}
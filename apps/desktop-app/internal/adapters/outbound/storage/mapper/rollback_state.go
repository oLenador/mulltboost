
package storage

import (
	"github.com/oLenador/mulltbost/internal/adapters/outbound/storage/models"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"gorm.io/datatypes"
)

func MapRollbackToDomain(s *storage.BoosterRollbackState) *entities.BoosterRollbackState {
	if s == nil {
		return nil
	}

	var backup map[string]interface{}
	if s.BackupData != nil {
		backup = map[string]interface{}(s.BackupData)
	} else {
		backup = make(map[string]interface{})
	}

	return &entities.BoosterRollbackState{
		ID:         s.ID,
		Applied:    s.Applied,
		AppliedAt:  s.AppliedAt,
		RevertedAt: s.RevertedAt,
		Version:    s.Version,
		BackupData: backup,
		Status:     entities.ExecutionStatus(s.Status),
		ErrorMsg:   s.ErrorMsg,
	}
}

func MapRollbackFromDomain(e *entities.BoosterRollbackState) *storage.BoosterRollbackState {
	if e == nil {
		return nil
	}

	var backup datatypes.JSONMap
	if e.BackupData != nil {
		backup = datatypes.JSONMap(e.BackupData)
	} else {
		backup = datatypes.JSONMap{}
	}

	return &storage.BoosterRollbackState{
		ID:         e.ID,
		Applied:    e.Applied,
		AppliedAt:  e.AppliedAt,
		RevertedAt: e.RevertedAt,
		Version:    e.Version,
		BackupData: backup,
		Status:     storage.ExecutionStatus(e.Status),
		ErrorMsg:   e.ErrorMsg,
	}
}

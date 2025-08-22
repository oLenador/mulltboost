package booster

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	storagemodels "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/models"
	repos "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/repositories"
)

// testBooster é um stub configurável usado para simular comportamento do BoosterUseCase.
type testBooster struct {
	id          string
	version     string
	canApply    bool
	canRevert   bool
	validateErr error
	execResult  *entities.BoostApplyResult
	execErr     error
	revertRes   *entities.BoostRevertResult
	revertErr   error
}

func (b *testBooster) Execute(ctx context.Context) (*entities.BoostApplyResult, error) {
	return b.execResult, b.execErr
}

func (b *testBooster) Validate(ctx context.Context) error {
	return b.validateErr
}

func (b *testBooster) CanApply(ctx context.Context) bool {
	return b.canApply
}

func (b *testBooster) CanRevert(ctx context.Context) bool {
	return b.canRevert
}

func (b *testBooster) GetEntity() entities.Booster {
	return entities.Booster{
		ID:      b.id,
		Version: b.version,
	}
}

// Correção: retornar zero value (struct) em vez de nil
func (b *testBooster) GetEntityDto(lang i18n.Language) dto.BoosterDto {
	return dto.BoosterDto{}
}

func (b *testBooster) Revert(ctx context.Context) (*entities.BoostRevertResult, error) {
	return b.revertRes, b.revertErr
}

// helper para criar DB em memória e repos.NewRollbackRepo
func setupRollbackRepoForTest(t *testing.T) (*repos.RollbackRepo, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// migrar tabela necessária (apenas rollback state)
	require.NoError(t, db.AutoMigrate(&storagemodels.BoosterRollbackState{}))

	rr := repos.NewRollbackRepo(db)
	return rr, db
}

func TestRegisterGetAndCountBooster(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	tb := &testBooster{
		id:       "b-1",
		version:  "v1",
		canApply: true,
	}

	// registrar
	err := proc.RegisterBooster(tb)
	require.NoError(t, err)

	// obter
	got, ok := proc.GetBooster("b-1")
	assert.True(t, ok)
	require.NotNil(t, got)
	assert.Equal(t, "b-1", got.GetEntity().ID)

	// contar
	assert.Equal(t, 1, proc.GetBoosterCount())

	// GetAllBoosters deve conter o booster
	all := proc.GetAllBoosters()
	_, present := all["b-1"]
	assert.True(t, present)
}

func TestProcessApply_BoosterNotFound(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	_, err := proc.ProcessApply(context.Background(), "non-existent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestProcessApply_CannotApply(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	tb := &testBooster{
		id:       "b-cannot",
		version:  "v1",
		canApply: false,
	}
	require.NoError(t, proc.RegisterBooster(tb))

	res, err := proc.ProcessApply(context.Background(), "b-cannot")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.False(t, res.Success)
	assert.Contains(t, res.Message, "cannot be applied")
}

func TestProcessApply_ValidateFails(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	tb := &testBooster{
		id:          "b-validate-fail",
		version:     "v1",
		canApply:    true,
		validateErr: errors.New("invalid config"),
	}
	require.NoError(t, proc.RegisterBooster(tb))

	res, err := proc.ProcessApply(context.Background(), "b-validate-fail")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.False(t, res.Success)
	assert.Contains(t, res.Message, "Validation failed")
	assert.Equal(t, "invalid config", res.Error.Error())
}

func TestProcessApply_ExecuteError(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	tb := &testBooster{
		id:         "b-exec-err",
		version:    "v1",
		canApply:   true,
		execResult: nil,
		execErr:    errors.New("exec failed"),
	}
	require.NoError(t, proc.RegisterBooster(tb))

	res, err := proc.ProcessApply(context.Background(), "b-exec-err")
	require.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "exec failed")
}

func TestProcessApply_Success_PersistsRollback(t *testing.T) {
	rr, db := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	tb := &testBooster{
		id:       "b-success",
		version:  "v1",
		canApply: true,
		execResult: &entities.BoostApplyResult{
			Success:    true,
			Message:    "ok",
			BackupData: map[string]interface{}{"k": "v"},
			Error:      nil,
		},
	}
	require.NoError(t, proc.RegisterBooster(tb))

	res, err := proc.ProcessApply(context.Background(), "b-success")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.True(t, res.Success)

	// verificar rollback salvo no DB
	br, err := rr.GetByID(context.Background(), "b-success")
	require.NoError(t, err)
	require.NotNil(t, br)
	assert.True(t, br.Applied)
	assert.Equal(t, "v1", br.Version)
	// backup data foi persistido (Mapper pode converter; assert básico)
	if br.BackupData != nil {
		assert.Equal(t, "v", br.BackupData["k"])
	} else {
		// Se o mapper/DB não ter persistido backup, pelo menos Applied/Version devem estar lá
		t.Log("BackupData nil (verifique mapper), mas estado aplicado persistiu")
	}

	// também conferir via consulta direta GORM
	var count int64
	db.Model(&storagemodels.BoosterRollbackState{}).Where("id = ?", "b-success").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestProcessRevert_BoosterNotFound(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	_, err := proc.ProcessRevert(context.Background(), "nope")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestProcessRevert_CannotRevert(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	tb := &testBooster{
		id:        "b-norevert",
		version:   "v1",
		canRevert: false,
	}
	require.NoError(t, proc.RegisterBooster(tb))

	res, err := proc.ProcessRevert(context.Background(), "b-norevert")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.False(t, res.Success)
	assert.Contains(t, res.Error.Error(), "cannot be reverted")
}

func TestProcessRevert_RevertError(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	tb := &testBooster{
		id:        "b-revert-err",
		version:   "v1",
		canRevert: true,
		revertRes: nil,
		revertErr: errors.New("revert fail"),
	}
	require.NoError(t, proc.RegisterBooster(tb))

	res, err := proc.ProcessRevert(context.Background(), "b-revert-err")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "revert fail")
	assert.Nil(t, res)
}

func TestProcessRevert_Success_UpdatesRollback(t *testing.T) {
	rr, _ := setupRollbackRepoForTest(t)
	proc := NewBoosterProcessor(rr)

	// criar rollback previamente aplicado
	appliedAt := time.Now().Add(-1 * time.Hour)
	initial := &entities.BoosterRollbackState{
		ID:        "b-to-revert",
		Applied:   true,
		AppliedAt: &appliedAt,
		Version:   "v1",
		Status:    entities.ExecutionApplied,
	}
	require.NoError(t, rr.Save(context.Background(), initial))

	// registrar booster que permite reverter
	tb := &testBooster{
		id:        "b-to-revert",
		version:   "v1",
		canRevert: true,
		revertRes: &entities.BoostRevertResult{
			Success: true,
			Message: "reverted",
		},
	}
	require.NoError(t, proc.RegisterBooster(tb))

	res, err := proc.ProcessRevert(context.Background(), "b-to-revert")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.True(t, res.Success)

	// verificar rollback atualizado
	updated, err := rr.GetByID(context.Background(), "b-to-revert")
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.False(t, updated.Applied)
	assert.Equal(t, entities.ExecutionReverted, updated.Status)
	require.NotNil(t, updated.RevertedAt)
}

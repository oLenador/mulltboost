package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	storage "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/models"
	repo "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/repositories"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrar tabelas
	err = db.AutoMigrate(&storage.BoostActivationState{}, &storage.BoosterRollbackState{})
	require.NoError(t, err)

	return db
}

func TestBoostConfigRepository_SyncWithAvailableBoosts(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Criar boosts de teste
	boosts := []entities.Booster{
		{ID: "test-boost-1", Version: "1.0.0"},
		{ID: "test-boost-2", Version: "1.1.0"},
	}

	repository := repo.NewBoostConfigRepository(db, boosts)

	t.Run("Deve criar novos estados para boosts", func(t *testing.T) {
		err := repository.SyncWithAvailableBoosts(ctx, boosts)
		require.NoError(t, err)

		// Verificar se os estados foram criados
		var count int64
		db.Model(&storage.BoostActivationState{}).Count(&count)
		assert.Equal(t, int64(2), count)

		// Verificar estado específico
		state, err := repository.GetBoostState(ctx, "test-boost-1")
		require.NoError(t, err)
		assert.False(t, state.IsApplied)
		assert.Equal(t, entities.StatusInactive, state.Status)
	})

	t.Run("Deve ativar e desativar boost", func(t *testing.T) {
		// Ativar boost
		err := repository.ActivateBoost(ctx, "test-boost-1")
		require.NoError(t, err)

		// Verificar se foi ativado
		assert.True(t, repository.IsBoostActive(ctx, "test-boost-1"))
		state, err := repository.GetBoostState(ctx, "test-boost-1")
		require.NoError(t, err)
		assert.True(t, state.IsApplied)
		assert.Equal(t, entities.StatusActive, state.Status)
		assert.NotNil(t, state.AppliedAt)

		// Desativar boost
		err = repository.DeactivateBoost(ctx, "test-boost-1")
		require.NoError(t, err)

		// Verificar se foi desativado
		assert.False(t, repository.IsBoostActive(ctx, "test-boost-1"))
		state, err = repository.GetBoostState(ctx, "test-boost-1")
		require.NoError(t, err)
		assert.False(t, state.IsApplied)
		assert.Equal(t, entities.StatusInactive, state.Status)
		assert.NotNil(t, state.RevertedAt)
	})

	t.Run("Deve marcar boosts removidos como obsoletos", func(t *testing.T) {
		// Sincronizar apenas com um boost
		singleBoost := []entities.Booster{boosts[0]}
		err := repository.SyncWithAvailableBoosts(ctx, singleBoost)
		require.NoError(t, err)

		// Verificar se o segundo boost foi marcado como obsoleto
		var obsoleteState storage.BoostActivationState
		err = db.Where("id = ? AND status = ?", "test-boost-2", entities.StatusObsolete).First(&obsoleteState).Error
		require.NoError(t, err)
		assert.False(t, obsoleteState.IsApplied)
	})

	t.Run("Deve limpar estados obsoletos antigos", func(t *testing.T) {
		// Criar estado obsoleto antigo
		oldState := storage.BoostActivationState{
			ID:        "old-boost",
			IsApplied: false,
			Status:    entities.StatusObsolete,
			Version:   "1.0.0",
			CreatedAt: time.Now().AddDate(0, 0, -10), // 10 dias atrás
			UpdatedAt: time.Now().AddDate(0, 0, -10),
		}
		db.Create(&oldState)

		// Executar limpeza
		err := repository.CleanupObsoleteStates(ctx)
		require.NoError(t, err)

		// Verificar se foi removido
		var count int64
		db.Model(&storage.BoostActivationState{}).Where("id = ?", "old-boost").Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestBoostConfigRepository_ErrorHandling(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	boosts := []entities.Booster{{ID: "test-boost", Version: "1.0.0"}}
	repository := repo.NewBoostConfigRepository(db, boosts)

	t.Run("Deve retornar erro ao ativar boost inexistente", func(t *testing.T) {
		err := repository.ActivateBoost(ctx, "non-existent-boost")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não encontrado")
	})

	t.Run("Deve retornar erro ao ativar boost já ativo", func(t *testing.T) {
		// Ativar primeiro
		err := repository.ActivateBoost(ctx, "test-boost")
		require.NoError(t, err)

		// Tentar ativar novamente
		err = repository.ActivateBoost(ctx, "test-boost")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "já está ativo")
	})

	t.Run("Deve retornar erro ao desativar boost já inativo", func(t *testing.T) {
		// Primeiro desativa
		err := repository.DeactivateBoost(ctx, "test-boost")
		require.NoError(t, err)

		// Tentar desativar novamente
		err = repository.DeactivateBoost(ctx, "test-boost")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "já está inativo")
	})
}


func setupTestDBForMoreTests(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrar tabelas (inclui BoostActivationState e BoosterRollbackState)
	err = db.AutoMigrate(&storage.BoostActivationState{}, &storage.BoosterRollbackState{})
	require.NoError(t, err)

	return db
}

// 1) Quando existir um BoosterRollbackState com Applied = true e RevertedAt = nil,
// a ativação do boost deve ser bloqueada.
func TestActivateBlockedByRollbackState(t *testing.T) {
	db := setupTestDBForMoreTests(t)
	ctx := context.Background()

	// criar boost nos "boosts" passados ao repo para que ele registre o estado
	boosts := []entities.Booster{{ID: "blocked-boost", Version: "1.0.0"}}
	repository := repo.NewBoostConfigRepository(db, boosts)

	// criar rollback state que bloqueia re-aplicação
	rollback := storage.BoosterRollbackState{
		ID:        "blocked-boost",
		Applied:   true,
		AppliedAt: ptrTime(time.Now().Add(-24 * time.Hour)),
		// RevertedAt == nil -> bloqueia
		Version: "1.0.0",
	}
	require.NoError(t, db.Create(&rollback).Error)

	// tentar ativar: deve falhar por rollback anterior não revertido
	err := repository.ActivateBoost(ctx, "blocked-boost")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "não é possível ativar")
}

// 2) GetBoostState: se não estiver em cache, deve buscar no DB
func TestGetBoostState_CacheMissReadsDB(t *testing.T) {
	db := setupTestDBForMoreTests(t)
	ctx := context.Background()

	// criar um estado apenas no DB
	now := time.Now()
	dbState := storage.BoostActivationState{
		ID:        "db-only",
		IsApplied: false,
		Version:   "v1",
		Status:    entities.StatusInactive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.NoError(t, db.Create(&dbState).Error)

	// criar repo sem passar esse boost (fica fora do cache)
	repository := repo.NewBoostConfigRepository(db, []entities.Booster{})

	// buscar diretamente
	got, err := repository.GetBoostState(ctx, "db-only")
	require.NoError(t, err)
	assert.Equal(t, "db-only", got.ID)
	assert.False(t, got.IsApplied)

	// agora que foi buscado, deve estar no cache (IsBoostActive usa cache)
	assert.False(t, repository.IsBoostActive(ctx, "db-only"))
}

// 3) CleanupObsoleteStates: só remove os obsoletos mais antigos que o limite
func TestCleanupObsoleteStates_RemovesOlderThanThreshold(t *testing.T) {
	db := setupTestDBForMoreTests(t)
	ctx := context.Background()

	// repo com nenhum boost
	repository := repo.NewBoostConfigRepository(db, []entities.Booster{})

	// criar 2 estados obsoletos: um velho, outro recente
	old := storage.BoostActivationState{
		ID:        "old-boost",
		IsApplied: false,
		Version:   "1.0.0",
		Status:    entities.StatusObsolete,
		CreatedAt: time.Now().AddDate(0, 0, -10),
		UpdatedAt: time.Now().AddDate(0, 0, -10),
	}
	recent := storage.BoostActivationState{
		ID:        "recent-boost",
		IsApplied: false,
		Version:   "1.0.1",
		Status:    entities.StatusObsolete,
		CreatedAt: time.Now().AddDate(0, 0, -1),
		UpdatedAt: time.Now().AddDate(0, 0, -1),
	}
	require.NoError(t, db.Create(&old).Error)
	require.NoError(t, db.Create(&recent).Error)

	// executar limpeza
	require.NoError(t, repository.CleanupObsoleteStates(ctx))

	// verificar que o velho foi removido e o recente não
	var cntOld int64
	db.Model(&storage.BoostActivationState{}).Where("id = ?", "old-boost").Count(&cntOld)
	assert.Equal(t, int64(0), cntOld)

	var cntRecent int64
	db.Model(&storage.BoostActivationState{}).Where("id = ?", "recent-boost").Count(&cntRecent)
	assert.Equal(t, int64(1), cntRecent)
}

// 4) Ativar boost inexistente deve retornar erro
func TestActivate_NonExistentBoost(t *testing.T) {
	db := setupTestDBForMoreTests(t)
	ctx := context.Background()
	repository := repo.NewBoostConfigRepository(db, []entities.Booster{})

	err := repository.ActivateBoost(ctx, "does-not-exist")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "não encontrado")
}

// 5) SetCurrentVersion não deve provocar pânico e deve permitir re-sync sem erro
func TestSetCurrentVersionAndSync(t *testing.T) {
	db := setupTestDBForMoreTests(t)
	ctx := context.Background()

	boosts := []entities.Booster{{ID: "b1", Version: "1.0.0"}}
	repository := repo.NewBoostConfigRepository(db, boosts)

	// trocar versão atual do repo e fazer um sync com um boost de versão diferente
	repository.SetCurrentVersion("2.0.0")

	// simular que a versão do código mudou (novo version string)
	newBoosts := []entities.Booster{{ID: "b1", Version: "2.0.0"}}
	require.NoError(t, repository.SyncWithAvailableBoosts(ctx, newBoosts))

	// buscar estado e checar que a versão foi atualizada no DB/cache
	st, err := repository.GetBoostState(ctx, "b1")
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", st.Version)
}
func ptrTime(t time.Time) *time.Time { return &t }
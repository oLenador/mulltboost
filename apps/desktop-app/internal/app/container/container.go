// internal/app/container/container.go
package container

import (
	"fmt"

	storage "github.com/oLenador/mulltbost/internal/adapters/outbound/storage" 
	models  "github.com/oLenador/mulltbost/internal/adapters/outbound/storage/models" 
	repos   "github.com/oLenador/mulltbost/internal/adapters/outbound/storage/repositories" 

	"github.com/oLenador/mulltbost/internal/adapters/outbound/system"
	"github.com/oLenador/mulltbost/internal/boosters/connection"

	"github.com/oLenador/mulltbost/internal/core/domain/services/booster"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/domain/services/monitoring"
	"github.com/oLenador/mulltbost/internal/core/domain/services/sysinfo"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

type Container struct {
	// Services
	BoosterService     inbound.BoosterService
	MonitoringService  inbound.MonitoringService
	SystemInfoService  inbound.SystemInfoService
	I18nService        *i18n.Service

	// Repositories
	metricsRepo    *system.MetricsRepository
	systemInfoRepo *system.InfoRepository
}

func NewContainer() (*Container, error) {
	// Repositories
	db, err := storage.NewDB()
	if err != nil {
		fmt.Printf("error on open db: %v", err)
		return nil, err
	}

	// 2) Executa migrações para criar as tabelas necessárias
	if err := storage.AutoMigrateModels(db, &models.BoosterRollbackState{}, &models.AppliedBoost{}); err != nil {
		fmt.Printf("automigrate : %v", err)
		return nil, err	
	}

	rollbackRepo := repos.NewRollbackRepo(db)
	appliedRepo := repos.NewAppliedRepo(db)
	
	metricsRepo := system.NewMetricsRepository()
	systemInfoRepo := system.NewInfoRepository()

	// Services
	i18nService := i18n.NewService()
	boosterService := booster.NewService(rollbackRepo, appliedRepo)

	monitoringService := monitoring.NewService(metricsRepo)
	systemInfoService := systeminfo.NewService(systemInfoRepo)

	container := &Container{
		BoosterService:    boosterService,
		MonitoringService: monitoringService,
		SystemInfoService: systemInfoService,
		I18nService:       i18nService,

		metricsRepo:    metricsRepo,
		systemInfoRepo: systemInfoRepo,
	}

	// Inicializar plugins direto no construtor
	if err := container.initAllBoosts(); err != nil {
		return nil, fmt.Errorf("failed to initialize plugins: %w", err)
	}

	return container, nil
}

func (c *Container) initAllBoosts() error {
	loaders := map[string][]inbound.BoosterUseCase{
		"connection": connection.GetAllPlugins(),
	}

	for _, boostArray := range loaders {
		for _, booster := range boostArray {
			if err := c.BoosterService.RegisterBooster(booster); err != nil {
				return err
			}
		}
	}
	return nil
}

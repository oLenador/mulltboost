// internal/app/container/container.go
package container

import (
	"fmt"

	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/services/booster"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/domain/services/monitoring"
	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection"
	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/system"
	"github.com/wailsapp/wails/v3/pkg/application"

	boosterBase "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
	storage "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage"
	models "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/models"
	repos "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/storage/repositories"
)

type Container struct {
	// Services
	BoosterService    inbound.BoosterService
	MetricsService    inbound.MonitoringService
	SystemInfoService inbound.SystemInfoService
	I18nService       *i18n.Service
	// Repositories
}

func NewContainer(appService *application.App) (*Container, error) {

	db, err := storage.NewDB()
	if err != nil {
		fmt.Printf("error on open db: %v", err)
		return nil, err
	}

	if err := storage.AutoMigrateModels(db, &models.BoosterRollbackState{}, &models.BoostOperation{}); err != nil {
		fmt.Printf("automigrate : %v", err)
		return nil, err
	}

	rollbackRepo := repos.NewRollbackRepo(db)
	boostOperationsRepo := repos.NewBoostOperationsRepo(db)

	systemMetricsRepo := system.NewMetricsRepository()
	metricsService := monitoring.NewService(systemMetricsRepo)

	// Services
	i18nService := i18n.NewService()
	boosterService := booster.NewService(rollbackRepo, boostOperationsRepo, appService.Event)

	container := &Container{
		BoosterService: boosterService,
		MetricsService: metricsService,
		I18nService:    i18nService,
	}

	// Inicializar plugins direto no construtor
	if err := container.initAllBoosts(); err != nil {
		return nil, fmt.Errorf("failed to initialize plugins: %w", err)
	}

	return container, nil
}

func (c *Container) initAllBoosts() error {

	ps := boosterBase.GetPlatformServices()
	deps := inbound.NewExecutorDepServices(ps)

	loaders := map[string][]inbound.BoosterUseCase{
		"connection": connection.GetAllPlugins(deps),
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

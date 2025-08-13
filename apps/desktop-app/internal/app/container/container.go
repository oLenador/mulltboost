// internal/app/container/container.go
package container

import (
	"fmt"

	"github.com/oLenador/mulltbost/internal/adapters/outbound/storage"
	"github.com/oLenador/mulltbost/internal/adapters/outbound/system"
	"github.com/oLenador/mulltbost/internal/boosters/connection"
	"github.com/oLenador/mulltbost/internal/boosters/flusher"
	"github.com/oLenador/mulltbost/internal/config"

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
	stateRepo      *storage.BoosterStateRepository
	metricsRepo    *system.MetricsRepository
	systemInfoRepo *system.InfoRepository
}

func NewContainer() (*Container, error) {
	// Repositories
	stateRepo, _ := storage.NewBoosterStateRepository(config.AppName)
	metricsRepo := system.NewMetricsRepository()
	systemInfoRepo := system.NewInfoRepository()

	// Services
	i18nService := i18n.NewService()
	boosterService := booster.NewService(stateRepo)

	monitoringService := monitoring.NewService(metricsRepo)
	systemInfoService := systeminfo.NewService(systemInfoRepo)

	container := &Container{
		BoosterService:    boosterService,
		MonitoringService: monitoringService,
		SystemInfoService: systemInfoService,
		I18nService:       i18nService,

		stateRepo:      stateRepo,
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
		"flusher":    flusher.GetAllPlugins(),
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

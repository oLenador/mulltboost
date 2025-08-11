package container

import (
    "github.com/oLenador/mulltbost/internal/adapters/outbound/storage"
    "github.com/oLenador/mulltbost/internal/adapters/outbound/system"
    "github.com/oLenador/mulltbost/internal/core/ports/inbound"
    "github.com/oLenador/mulltbost/internal/core/domain/usecases/optimization"
    "github.com/oLenador/mulltbost/internal/core/domain/usecases/monitoring"
    "github.com/oLenador/mulltbost/internal/core/domain/usecases/sysinfo"
    
    // Plugin imports
    "github.com/oLenador/mulltbost/internal/plugins/precision"
    "github.com/oLenador/mulltbost/internal/plugins/performance"
    "github.com/oLenador/mulltbost/internal/plugins/network"
    pluginSystem "github.com/oLenador/mulltbost/internal/plugins/system"
)


type Container struct {
    // Services
    OptimizationService inbound.OptimizationService
    MonitoringService   inbound.MonitoringService
    SystemInfoService   inbound.SystemInfoService
    
    // Repositories
    stateRepo       *storage.OptimizationStateRepository
    metricsRepo     *system.MetricsRepository
    systemInfoRepo  *system.InfoRepository
}

func NewContainer() *Container {
    // Repositories
    stateRepo := storage.NewOptimizationStateRepository()
    metricsRepo := system.NewMetricsRepository()
    systemInfoRepo := system.NewInfoRepository()
    
    // Services
    optimizationService := optimization.NewService(stateRepo)
    monitoringService := monitoring.NewService(metricsRepo)
    systemInfoService := systeminfo.NewService(systemInfoRepo)
    
    return &Container{
        OptimizationService: optimizationService,
        MonitoringService:   monitoringService,
        SystemInfoService:   systemInfoService,
        stateRepo:          stateRepo,
        metricsRepo:        metricsRepo,
        systemInfoRepo:     systemInfoRepo,
    }
}

func (c *Container) RegisterPlugins() error {
    optimizationService := c.OptimizationService.(*optimization.Service)
    
    // Precision plugins
    precisionPlugins := precision.GetAllPlugins()
    for _, plugin := range precisionPlugins {
        if err := optimizationService.RegisterPlugin(plugin); err != nil {
            return err
        }
    }
    
    // Performance plugins
    performancePlugins := performance.GetAllPlugins()
    for _, plugin := range performancePlugins {
        if err := optimizationService.RegisterPlugin(plugin); err != nil {
            return err
        }
    }
    
    // Network plugins
    networkPlugins := network.GetAllPlugins()
    for _, plugin := range networkPlugins {
        if err := optimizationService.RegisterPlugin(plugin); err != nil {
            return err
        }
    }
    
    // System plugins
    systemPlugins := pluginSystem.GetAllPlugins()
    for _, plugin := range systemPlugins {
        if err := optimizationService.RegisterPlugin(plugin); err != nil {
            return err
        }
    }
    
    return nil
}
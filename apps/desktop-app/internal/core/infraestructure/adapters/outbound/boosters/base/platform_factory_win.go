//go:build windows

package booster

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	outbound "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/system"
	winOutbound  "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/windows"
	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/windows"
)

func init() {
	DefaultFactory = &WindowsServicesFactory{}
}

type WindowsServicesFactory struct{}


func (f *WindowsServicesFactory) CreatePlatformServices() inbound.PlatformServices {
	netWorkService, _ := system.NewWindowsNetworkAdapterService()
	elevationService, _ := system.NewWindowsElevationService()
	powerService, _ := system.NewWindowsPowerManagementService()
	gpuInfoService := system.NewGPUInfoRepository()
	gpuService := system.NewGPUOptimizationService(gpuInfoService)

	return &WindowsPlatformServices{
		services: &inbound.ExecutorDepServices{
			TcpService:       system.NewTCPOptimizationService(),
			RegistryService:  system.NewRegistryService(),
			SystemService:    system.NewWindowsSystemService(),
			ElevationService: elevationService,
			NetworkService:   netWorkService,
			MemoryService:    system.NewWindowsMemoryManagementService(),
			GpuInfoService:   gpuInfoService,
			GpuService:       gpuService,
			PowerService:     powerService,
			ProcessService:   system.NewProcessPriorityService(),
			CacheService:     system.NewCacheManagementService(),
			FileService:      system.NewWindowsFileSystemService(),
			WindowsSvcMgr:    windows.NewWinServiceManagerService(),
		},
	}
}

func (f *WindowsServicesFactory) GetPlatform() string {
	return "windows"
}

type WindowsPlatformServices struct {
	services *inbound.ExecutorDepServices
}

func (w *WindowsPlatformServices) GetGpuInfoService() outbound.GPUInfoRepository {
	return w.services.GpuInfoService
}

func (w *WindowsPlatformServices) GetCacheService() outbound.CacheManagementService {
	return w.services.CacheService
}

func (w *WindowsPlatformServices) GetElevationService() outbound.ElevationService {
	return w.services.ElevationService
}

func (w *WindowsPlatformServices) GetFileService() outbound.FileSystemService {
	return w.services.FileService
}

func (w *WindowsPlatformServices) GetGpuService() outbound.GPUOptimizationService {
	return w.services.GpuService
}

func (w *WindowsPlatformServices) GetMemoryService() outbound.MemoryManagementService {
	return w.services.MemoryService
}

func (w *WindowsPlatformServices) GetNetworkService() outbound.NetworkAdapterService {
	return w.services.NetworkService
}

func (w *WindowsPlatformServices) GetPowerService() outbound.PowerManagementService {
	return w.services.PowerService
}

func (w *WindowsPlatformServices) GetProcessService() outbound.ProcessPriorityService {
	return w.services.ProcessService
}

func (w *WindowsPlatformServices) GetRegistryService() outbound.RegistryService {
	return w.services.RegistryService
}

func (w *WindowsPlatformServices) GetSystemService() outbound.SystemAPIService {
	return w.services.SystemService
}

func (w *WindowsPlatformServices) GetTcpService() outbound.TCPOptimizationService {
	return w.services.TcpService
}

func (w *WindowsPlatformServices) GetWindowsSvcMgr() winOutbound.WinServiceManagerService {
	return w.services.WindowsSvcMgr
}

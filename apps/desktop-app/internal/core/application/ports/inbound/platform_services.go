package inbound

import (
	outbound "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"github.com/oLenador/mulltbost/internal/core/application/ports/outbound/windows"
)

type PlatformServices interface {
	GetTcpService() outbound.TCPOptimizationService
	GetRegistryService() outbound.RegistryService
	GetSystemService() outbound.SystemAPIService
	GetElevationService() outbound.ElevationService
	GetNetworkService() outbound.NetworkAdapterService
	GetMemoryService() outbound.MemoryManagementService
	GetGpuInfoService() outbound.GPUInfoRepository
	GetGpuService() outbound.GPUOptimizationService
	GetPowerService() outbound.PowerManagementService
	GetProcessService() outbound.ProcessPriorityService
	GetCacheService() outbound.CacheManagementService
	GetFileService() outbound.FileSystemService
	GetWindowsSvcMgr() windows.WinServiceManagerService
}

type ExecutorDepServices struct {
	TcpService       outbound.TCPOptimizationService
	RegistryService  outbound.RegistryService
	SystemService    outbound.SystemAPIService
	ElevationService outbound.ElevationService
	NetworkService   outbound.NetworkAdapterService
	MemoryService    outbound.MemoryManagementService
	GpuService       outbound.GPUOptimizationService
	GpuInfoService   outbound.GPUInfoRepository
	PowerService     outbound.PowerManagementService
	ProcessService   outbound.ProcessPriorityService
	CacheService     outbound.CacheManagementService
	FileService      outbound.FileSystemService
	WindowsSvcMgr    windows.WinServiceManagerService
}

// Adapter para converter PlatformServices em ExecutorDepServices
func NewExecutorDepServices(ps PlatformServices) *ExecutorDepServices {
	return &ExecutorDepServices{
		TcpService:       ps.GetTcpService(),
		RegistryService:  ps.GetRegistryService(),
		SystemService:    ps.GetSystemService(),
		ElevationService: ps.GetElevationService(),
		NetworkService:   ps.GetNetworkService(),
		MemoryService:    ps.GetMemoryService(),
		GpuInfoService:   ps.GetGpuInfoService(),
		GpuService:       ps.GetGpuService(),
		PowerService:     ps.GetPowerService(),
		ProcessService:   ps.GetProcessService(),
		CacheService:     ps.GetCacheService(),
		FileService:      ps.GetFileService(),
		WindowsSvcMgr:    ps.GetWindowsSvcMgr(),
	}
}

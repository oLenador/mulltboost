//go:build linux
package booster

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	outbound "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
    winOutbound  "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/windows"

)

func init() {
	DefaultFactory = &LinuxServicesFactory{}
}

type LinuxServicesFactory struct{}


func (f *LinuxServicesFactory) CreatePlatformServices() inbound.PlatformServices {
	return &LinuxPlatformServices{
		services: &inbound.ExecutorDepServices{},
	}
}

func (f *LinuxServicesFactory) GetPlatform() string {
	return "linux"
}

type LinuxPlatformServices struct {
	services *inbound.ExecutorDepServices
}

func (l *LinuxPlatformServices) GetTcpService() outbound.TCPOptimizationService {
	return l.services.TcpService // será nil
}

func (l *LinuxPlatformServices) GetRegistryService() outbound.RegistryService {
	return l.services.RegistryService // será nil
}

func (l *LinuxPlatformServices) GetSystemService() outbound.SystemAPIService {
	return l.services.SystemService // será nil
}

func (l *LinuxPlatformServices) GetElevationService() outbound.ElevationService {
	return l.services.ElevationService // será nil
}

func (l *LinuxPlatformServices) GetNetworkService() outbound.NetworkAdapterService {
	return l.services.NetworkService
}

func (l *LinuxPlatformServices) GetMemoryService() outbound.MemoryManagementService {
	return l.services.MemoryService // será nil
}

func (l *LinuxPlatformServices) GetGpuInfoService() outbound.GPUInfoRepository {
	return l.services.GpuInfoService // será nil
}

func (l *LinuxPlatformServices) GetGpuService() outbound.GPUOptimizationService {
	return l.services.GpuService // será nil
}

func (l *LinuxPlatformServices) GetPowerService() outbound.PowerManagementService {
	return l.services.PowerService // será nil
}

func (l *LinuxPlatformServices) GetProcessService() outbound.ProcessPriorityService {
	return l.services.ProcessService
}

func (l *LinuxPlatformServices) GetCacheService() outbound.CacheManagementService {
	return l.services.CacheService // será nil
}

func (l *LinuxPlatformServices) GetFileService() outbound.FileSystemService {
	return l.services.FileService
}

func (l *LinuxPlatformServices) GetWindowsSvcMgr() winOutbound.WinServiceManagerService {
	return l.services.WindowsSvcMgr // será nil
}

func (l *LinuxPlatformServices) GetPlatform() string {
	return "linux"
}

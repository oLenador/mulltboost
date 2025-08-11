package system

import (
    "context"
    "runtime"

    "github.com/shirou/gopsutil/v4/cpu"
    "github.com/shirou/gopsutil/v4/host"
    "github.com/shirou/gopsutil/v4/mem"
    "github.com/shirou/gopsutil/v4/disk"
    "github.com/shirou/gopsutil/v4/net"

    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type InfoRepository struct{}

func NewInfoRepository() *InfoRepository {
    return &InfoRepository{}
}

func (r *InfoRepository) GetOSInfo(ctx context.Context) (*entities.OSInfo, error) {
    info, err := host.Info()
    if err != nil {
        return nil, err
    }

    return &entities.OSInfo{
        Name:         info.OS,       // "windows", "linux", "darwin"
        Version:      info.PlatformVersion,
        Architecture: runtime.GOARCH,
        BuildNumber:  info.KernelVersion,
    }, nil
}

func (r *InfoRepository) GetCPUInfo(ctx context.Context) (*entities.CPUInfo, error) {
    infos, err := cpu.Info()
    if err != nil || len(infos) == 0 {
        return nil, err
    }

    c := infos[0]
    counts, _ := cpu.Counts(true) // threads lógicos
    cores, _ := cpu.Counts(false) // núcleos físicos

    return &entities.CPUInfo{
        Name:         c.ModelName,
        Manufacturer: c.VendorID,
        Cores:        cores,
        Threads:      counts,
        BaseFreq:     c.Mhz,
        Architecture: runtime.GOARCH,
    }, nil
}

func (r *InfoRepository) GetMemoryInfo(ctx context.Context) (*entities.MemoryInfo, error) {
    vm, err := mem.VirtualMemory()
    if err != nil {
        return nil, err
    }

    return &entities.MemoryInfo{
        TotalRAM: vm.Total,
        // Slots, tipo e velocidade precisam de APIs específicas do SO
    }, nil
}

func (r *InfoRepository) GetStorageInfo(ctx context.Context) ([]entities.StorageInfo, error) {
    parts, err := disk.Partitions(true)
    if err != nil {
        return nil, err
    }

    var result []entities.StorageInfo
    for _, p := range parts {
        usage, _ := disk.Usage(p.Mountpoint)
        result = append(result, entities.StorageInfo{
            Name: p.Mountpoint,
            Type: p.Fstype,
            Size: usage.Total,
        })
    }
    return result, nil
}

func (r *InfoRepository) GetNetworkInfo(ctx context.Context) ([]entities.NetworkInfo, error) {
    ifaces, err := net.Interfaces()
    if err != nil {
        return nil, err
    }

    var result []entities.NetworkInfo
    for _, iface := range ifaces {
        for _, addr := range iface.Addrs {
            result = append(result, entities.NetworkInfo{
                Name:       iface.Name,
                MACAddress: iface.HardwareAddr,
                IPAddress:  addr.Addr,
            })
        }
    }
    return result, nil
}

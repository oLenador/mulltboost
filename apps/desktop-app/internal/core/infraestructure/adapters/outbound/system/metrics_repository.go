package system

import (
	"context"
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type MetricsRepository struct{}


func NewMetricsRepository() *MetricsRepository {
    return &MetricsRepository{}
}

func (r *MetricsRepository) GetCPUMetrics(ctx context.Context) (*entities.CPUMetrics, error) {

    usagePercents, err := cpu.Percent(200*time.Millisecond, true)
    if err != nil {
        return nil, err
    }

    // Frequências
    freqs, _ := cpu.Info()

    // Núcleos e threads
    cores, _ := cpu.Counts(false)
    threads, _ := cpu.Counts(true)

    coresMetrics := make([]entities.CoreMetric, len(usagePercents))
    for i, usage := range usagePercents {
        coresMetrics[i] = entities.CoreMetric{
            Index: i,
            Usage: usage,
            Frequency: func() float64 {
                if len(freqs) > i {
                    return freqs[i].Mhz
                }
                return 0
            }(),
        }
    }

    // Média de uso geral
    var totalUsage float64
    for _, u := range usagePercents {
        totalUsage += u
    }
    avgUsage := totalUsage / float64(len(usagePercents))

    return &entities.CPUMetrics{
        Usage:       avgUsage,
        CoreCount:   cores,
        ThreadCount: threads,
        Frequency:   freqs[0].Mhz,
        Cores:       coresMetrics,
    }, nil
}

func (r *MetricsRepository) GetMemoryMetrics(ctx context.Context) (*entities.MemoryMetrics, error) {
    vm, err := mem.VirtualMemory()
    if err != nil {
        return nil, err
    }

    return &entities.MemoryMetrics{
        Total:        vm.Total,
        Used:         vm.Used,
        Available:    vm.Available,
        UsagePercent: vm.UsedPercent,
        Cached:       vm.Cached,
        Buffers:      vm.Buffers,
    }, nil
}


func (r *MetricsRepository) GetDiskMetrics(ctx context.Context) (*entities.DiskMetrics, error) {
    parts, err := disk.Partitions(true)
    if err != nil {
        return nil, err
    }

    var drives []entities.DriveMetric
    for _, p := range parts {
        usage, err := disk.Usage(p.Mountpoint)
        if err != nil {
            continue
        }
        drives = append(drives, entities.DriveMetric{
            Name:         p.Mountpoint,
            Total:        usage.Total,
            Used:         usage.Used,
            Free:         usage.Free,
            UsagePercent: usage.UsedPercent,
        })
    }

    return &entities.DiskMetrics{Drives: drives}, nil
}

func (r *MetricsRepository) GetGPUMetrics(ctx context.Context) (*entities.GPUMetrics, error) {
    // gopsutil não fornece GPU — mock básico
    // Para NVIDIA: usar go-nvml ou chamar `nvidia-smi --query-gpu=...`
    return &entities.GPUMetrics{
        Name:  "Unknown GPU",
        Usage: 0.0, // sem suporte nativo aqui
        MemoryTotal: 0,
        MemoryUsed:  0,
    }, nil
}


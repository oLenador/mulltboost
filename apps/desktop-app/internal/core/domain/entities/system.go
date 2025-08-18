package entities

import "time"

type SystemMetrics struct {
    CPU         CPUMetrics
    Memory      MemoryMetrics
    GPU         GPUMetrics
    Network     NetworkMetrics
    Temperature TemperatureMetrics
    Disk        DiskMetrics
    Timestamp   time.Time
}

type CPUMetrics struct {
    Usage       float64
    CoreCount   int
    ThreadCount int
    Frequency   float64
    Temperature float64
    Cores       []CoreMetric
}

type CoreMetric struct {
    Index       int
    Usage       float64
    Frequency   float64
    Temperature float64
}

type MemoryMetrics struct {
    Total       uint64
    Used        uint64
    Available   uint64
    UsagePercent float64
    Cached      uint64
    Buffers     uint64
}

type GPUMetrics struct {
    Name            string
    Usage           float64
    MemoryUsed      uint64
    MemoryTotal     uint64
    Temperature     float64
    PowerDraw       float64
    ClockSpeed      int
    MemoryClockSpeed int
}

type NetworkInterface struct {
    Name      string
    BytesSent uint64
    BytesRecv uint64
    Speed     uint64
    IsUp      bool
}

type TemperatureMetrics struct {
    CPU        float64
    GPU        float64
    Motherboard float64
    Drives     []DriveTemp
}

type DriveTemp struct {
    Name        string
    Temperature float64
}

type DiskMetrics struct {
    Drives []DriveMetric
}

type DriveMetric struct {
    Name         string
    Total        uint64
    Used         uint64
    Free         uint64
    UsagePercent float64
    ReadSpeed    float64
    WriteSpeed   float64
}

type SystemInfo struct {
    OS          OSInfo
    CPU         CPUInfo
    Memory      MemoryInfo
    GPU         []GPUInfo
    Storage     []StorageInfo
    Network     []NetworkInfo
    Motherboard MotherboardInfo
}

type OSInfo struct {
    Name         string
    Version      string
    Architecture string
    BuildNumber  string
}

type CPUInfo struct {
    Name         string
    Manufacturer string
    Cores        int
    Threads      int
    BaseFreq     float64
    MaxFreq      float64
    Architecture string
}

type MemoryInfo struct {
    TotalRAM    uint64
    TotalSlots  int
    UsedSlots   int
    MemoryType  string
    Speed       int
}

type StorageInfo struct {
    Name       string
    Type       string
    Size       uint64
    Model      string
    Interface  string
}

type NetworkInfo struct {
    Name         string
    Type         string
    Speed        uint64
    MACAddress   string
    IPAddress    string
}

type MotherboardInfo struct {
    Manufacturer string
    Model        string
    BIOS         string
    Chipset      string
}
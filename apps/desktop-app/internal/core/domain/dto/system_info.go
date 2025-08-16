package dto

type SystemInfoDTO struct {
	OS               string            `json:"os" validate:"required"`
	Version          string            `json:"version"`
	Architecture     string            `json:"architecture"`
	Hostname         string            `json:"hostname"`
	TotalRAM         int64             `json:"total_ram"`
	AvailableRAM     int64             `json:"available_ram"`
	CPUModel         string            `json:"cpu_model"`
	CPUCores         int               `json:"cpu_cores"`
	CPUThreads       int               `json:"cpu_threads"`
	CPUUsage         float64           `json:"cpu_usage"`
	MemoryUsage      float64           `json:"memory_usage"`
	DiskUsage        map[string]float64 `json:"disk_usage"`
	GPUInfo          []GPUInfoDTO      `json:"gpu_info"`
	NetworkAdapters  []NetworkAdapterDTO `json:"network_adapters"`
	IsOptimized      bool              `json:"is_optimized"`
	LastOptimization *string           `json:"last_optimization,omitempty"`
}

type GPUInfoDTO struct {
	Name         string  `json:"name"`
	Vendor       string  `json:"vendor"`
	VRAMSize     int64   `json:"vram_size"`
	Usage        float64 `json:"usage"`
	Temperature  float64 `json:"temperature"`
	IsPrimary    bool    `json:"is_primary"`
}

type NetworkAdapterDTO struct {
	Name           string `json:"name"`
	IsConnected    bool   `json:"is_connected"`
	ConnectionType string `json:"connection_type"`
	Speed          int64  `json:"speed"`
	IPAddress      string `json:"ip_address"`
}

package entities

import "time"

type GPUInfo struct {
	ID               string    `json:"id" validate:"required"`
	Name             string    `json:"name"`
	Vendor           string    `json:"vendor"`
	DeviceID         string    `json:"device_id"`
	DriverVersion    string    `json:"driver_version"`
	VRAMSize         int64     `json:"vram_size"`
	VRAMUsed         int64     `json:"vram_used"`
	CoreClock        int       `json:"core_clock"`
	MemoryClock      int       `json:"memory_clock"`
	Temperature      float64   `json:"temperature"`
	Usage            float64   `json:"usage" validate:"min=0,max=100"`
	PowerUsage       float64   `json:"power_usage"`
	IsPrimary        bool      `json:"is_primary"`
	IsDiscrete       bool      `json:"is_discrete"`
	SupportsDirectX  string    `json:"supports_directx"`
	SupportsOpenGL   string    `json:"supports_opengl"`
	SupportsVulkan   bool      `json:"supports_vulkan"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

package entities

import "time"

type OptimizationProfile struct {
	ID                  string                 `json:"id" validate:"required"`
	Name                string                 `json:"name" validate:"required"`
	Description         string                 `json:"description"`
	Type                string                 `json:"type" validate:"oneof=gaming work custom balanced"`
	IsActive            bool                   `json:"is_active"`
	IsDefault           bool                   `json:"is_default"`
	Settings            map[string]interface{} `json:"settings"`
	CPUOptimization     bool                   `json:"cpu_optimization"`
	MemoryOptimization  bool                   `json:"memory_optimization"`
	GPUOptimization     bool                   `json:"gpu_optimization"`
	NetworkOptimization bool                   `json:"network_optimization"`
	AudioOptimization   bool                   `json:"audio_optimization"`
	PowerOptimization   bool                   `json:"power_optimization"`
	ServicesOptimization bool                  `json:"services_optimization"`
	RegistryOptimization bool                  `json:"registry_optimization"`
	Priority            int                    `json:"priority"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

package entities

import "time"

type WindowsService struct {
	ID              string    `json:"id" validate:"required"`
	ServiceName     string    `json:"service_name" validate:"required"`
	DisplayName     string    `json:"display_name"`
	Description     string    `json:"description"`
	Status          string    `json:"status" validate:"oneof=running stopped paused start_pending stop_pending"`
	StartupType     string    `json:"startup_type" validate:"oneof=automatic manual disabled delayed_automatic"`
	IsEssential     bool      `json:"is_essential"`
	IsOptimized     bool      `json:"is_optimized"`
	CanBeStopped    bool      `json:"can_be_stopped"`
	CanBePaused     bool      `json:"can_be_paused"`
	ProcessID       int       `json:"process_id"`
	MemoryUsage     int64     `json:"memory_usage"`
	Dependencies    []string  `json:"dependencies"`
	DependentServices []string `json:"dependent_services"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

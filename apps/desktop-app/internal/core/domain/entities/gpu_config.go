package entities

import "time"

type GPUConfiguration struct {
	ID                  string    `json:"id" validate:"required"`
	GPUID               string    `json:"gpu_id" validate:"required"`
	PowerMode           string    `json:"power_mode" validate:"oneof=balanced performance power_saver"`
	CoreClockOffset     int       `json:"core_clock_offset"`
	MemoryClockOffset   int       `json:"memory_clock_offset"`
	PowerLimitPercent   int       `json:"power_limit_percent" validate:"min=50,max=120"`
	FanCurveEnabled     bool      `json:"fan_curve_enabled"`
	OverclockingEnabled bool      `json:"overclocking_enabled"`
	IsOptimized         bool      `json:"is_optimized"`
	ProfileName         string    `json:"profile_name"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

package dto

type GPUConfigurationDTO struct {
	GPUID               string `json:"gpu_id" validate:"required"`
	PowerMode           string `json:"power_mode" validate:"oneof=balanced performance power_saver"`
	CoreClockOffset     int    `json:"core_clock_offset"`
	MemoryClockOffset   int    `json:"memory_clock_offset"`
	PowerLimitPercent   int    `json:"power_limit_percent" validate:"min=50,max=120"`
	FanCurveEnabled     bool   `json:"fan_curve_enabled"`
	OverclockingEnabled bool   `json:"overclocking_enabled"`
	ProfileName         string `json:"profile_name"`
}
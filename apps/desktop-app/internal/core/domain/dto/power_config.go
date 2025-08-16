package dto

type PowerConfigurationDTO struct {
	ProfileName             string `json:"profile_name" validate:"required"`
	CPUMinState             int    `json:"cpu_min_state" validate:"min=0,max=100"`
	CPUMaxState             int    `json:"cpu_max_state" validate:"min=0,max=100"`
	CPUPowerPolicy          string `json:"cpu_power_policy"`
	USBSelectiveSuspend     bool   `json:"usb_selective_suspend"`
	PowerThrottlingEnabled  bool   `json:"power_throttling_enabled"`
	HibernationEnabled      bool   `json:"hibernation_enabled"`
	SleepTimeout            int    `json:"sleep_timeout"`
	DisplayTimeout          int    `json:"display_timeout"`
	EnableFastStartup       bool   `json:"enable_fast_startup"`
}
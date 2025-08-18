package entities

import "time"

type PowerProfile struct {
	ID                    string    `json:"id" validate:"required"`
	Name                  string    `json:"name" validate:"required"`
	Description           string    `json:"description"`
	IsActive              bool      `json:"is_active"`
	IsCustom              bool      `json:"is_custom"`
	CPUMinState           int       `json:"cpu_min_state" validate:"min=0,max=100"`
	CPUMaxState           int       `json:"cpu_max_state" validate:"min=0,max=100"`
	CPUPowerPolicy        string    `json:"cpu_power_policy"`
	USBSelectiveSuspend   bool      `json:"usb_selective_suspend"`
	PowerThrottlingEnabled bool     `json:"power_throttling_enabled"`
	HibernationEnabled    bool      `json:"hibernation_enabled"`
	SleepTimeout          int       `json:"sleep_timeout"`
	DisplayTimeout        int       `json:"display_timeout"`
	IsOptimized           bool      `json:"is_optimized"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

package entities

import "time"

type DriverConfiguration struct {
	ID              string    `json:"id" validate:"required"`
	DeviceID        string    `json:"device_id" validate:"required"`
	DeviceName      string    `json:"device_name"`
	DriverVersion   string    `json:"driver_version"`
	DriverDate      time.Time `json:"driver_date"`
	DriverProvider  string    `json:"driver_provider"`
	IsOutdated      bool      `json:"is_outdated"`
	IsOptimized     bool      `json:"is_optimized"`
	UpdateAvailable bool      `json:"update_available"`
	LatestVersion   string    `json:"latest_version"`
	DeviceClass     string    `json:"device_class"`
	IsEnabled       bool      `json:"is_enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

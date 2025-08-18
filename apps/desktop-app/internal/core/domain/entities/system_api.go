package entities

import "time"

type SystemAPI struct {
	ID              string                 `json:"id" validate:"required"`
	APIName         string                 `json:"api_name" validate:"required"`
	DLLName         string                 `json:"dll_name"`
	IsLoaded        bool                   `json:"is_loaded"`
	Version         string                 `json:"version"`
	Parameters      map[string]interface{} `json:"parameters"`
	LastCallResult  interface{}            `json:"last_call_result"`
	LastCallTime    time.Time              `json:"last_call_time"`
	CallCount       int64                  `json:"call_count"`
	IsAvailable     bool                   `json:"is_available"`
	RequiresElevation bool                 `json:"requires_elevation"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

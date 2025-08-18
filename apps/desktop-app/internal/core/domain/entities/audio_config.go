package entities

import "time"

type AudioConfiguration struct {
	ID                string    `json:"id" validate:"required"`
	DeviceID          string    `json:"device_id" validate:"required"`
	DeviceName        string    `json:"device_name"`
	SampleRate        int       `json:"sample_rate" validate:"min=8000,max=192000"`
	BitDepth          int       `json:"bit_depth" validate:"oneof=16 24 32"`
	Latency           int       `json:"latency" validate:"min=1,max=1000"`
	BufferSize        int       `json:"buffer_size"`
	ExclusiveMode     bool      `json:"exclusive_mode"`
	EnhancementsEnabled bool    `json:"enhancements_enabled"`
	SpatialAudioEnabled bool    `json:"spatial_audio_enabled"`
	IsOptimized       bool      `json:"is_optimized"`
	ProfileName       string    `json:"profile_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

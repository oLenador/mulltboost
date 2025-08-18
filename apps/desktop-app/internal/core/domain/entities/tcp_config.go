package entities

import "time"

type TCPConfiguration struct {
	ID                    string    `json:"id" validate:"required"`
	WindowSize            int       `json:"window_size"`
	CongestionControl     string    `json:"congestion_control"`
	NagleAlgorithmEnabled bool      `json:"nagle_algorithm_enabled"`
	DelayedAckEnabled     bool      `json:"delayed_ack_enabled"`
	WindowScalingEnabled  bool      `json:"window_scaling_enabled"`
	SACKEnabled           bool      `json:"sack_enabled"`
	TimestampsEnabled     bool      `json:"timestamps_enabled"`
	IsOptimized           bool      `json:"is_optimized"`
	ProfileName           string    `json:"profile_name"`
	OptimizationType      string    `json:"optimization_type" validate:"oneof=gaming streaming general"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

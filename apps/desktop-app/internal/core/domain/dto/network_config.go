package dto

type NetworkConfigurationDTO struct {
	AdapterID           string   `json:"adapter_id" validate:"required"`
	Priority            int      `json:"priority"`
	DNSServers          []string `json:"dns_servers"`
	DisableThrottling   bool     `json:"disable_throttling"`
	OptimizeTCP         bool     `json:"optimize_tcp"`
	TCPConfiguration    TCPConfigurationDTO `json:"tcp_configuration"`
	QoSEnabled          bool     `json:"qos_enabled"`
	BandwidthLimit      int64    `json:"bandwidth_limit,omitempty"`
}

type TCPConfigurationDTO struct {
	WindowSize            int    `json:"window_size"`
	CongestionControl     string `json:"congestion_control"`
	NagleAlgorithmEnabled bool   `json:"nagle_algorithm_enabled"`
	DelayedAckEnabled     bool   `json:"delayed_ack_enabled"`
	WindowScalingEnabled  bool   `json:"window_scaling_enabled"`
	SACKEnabled           bool   `json:"sack_enabled"`
	TimestampsEnabled     bool   `json:"timestamps_enabled"`
	OptimizationType      string `json:"optimization_type" validate:"oneof=gaming streaming general"`
}
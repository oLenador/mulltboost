package dto

type ProcessManagementDTO struct {
	ProcessName       string   `json:"process_name" validate:"required"`
	ProcessID         int      `json:"process_id,omitempty"`
	Priority          int      `json:"priority"`
	AffinityMask      uint64   `json:"affinity_mask"`
	Action            string   `json:"action" validate:"oneof=optimize kill suspend resume"`
	CPULimit          float64  `json:"cpu_limit,omitempty"`
	MemoryLimit       int64    `json:"memory_limit,omitempty"`
	SetGameMode       bool     `json:"set_game_mode"`
	AllowedCPUCores   []int    `json:"allowed_cpu_cores,omitempty"`
}

type ProcessListDTO struct {
	Processes      []ProcessInfoDTO `json:"processes"`
	TotalProcesses int              `json:"total_processes"`
	SystemLoad     float64          `json:"system_load"`
	HighCPUProcesses []ProcessInfoDTO `json:"high_cpu_processes"`
	HighMemoryProcesses []ProcessInfoDTO `json:"high_memory_processes"`
}

type ProcessInfoDTO struct {
	ProcessID       int     `json:"process_id"`
	Name            string  `json:"name"`
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryUsage     int64   `json:"memory_usage"`
	ThreadCount     int     `json:"thread_count"`
	Priority        int     `json:"priority"`
	IsGameProcess   bool    `json:"is_game_process"`
	IsSystemProcess bool    `json:"is_system_process"`
	IsOptimized     bool    `json:"is_optimized"`
}
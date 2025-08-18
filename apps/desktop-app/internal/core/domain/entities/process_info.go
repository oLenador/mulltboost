package entities

import "time"

type ProcessInfo struct {
	PID            int  
	PPID           int
	ID              string    
	ProcessID       int       
	Name            string    
	Path            string    
	Priority        int       
	CPUUsage        float64   
	MemoryUsage     uint64    
	ThreadCount     int       
	HandleCount     int       
	IsOptimized     bool      
	AffinityMask    uint64    
	IsGameProcess   bool      
	IsSystemProcess bool      
	StartTime       time.Time 
	CreatedAt       time.Time 
	UpdatedAt       time.Time 
}

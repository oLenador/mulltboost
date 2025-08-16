package events

import "time"

const (
	SYSTEM_ERROR        string = "system_error"
	OPTIMIZATION_ERROR  string = "optimization_error"
	CONFIGURATION_ERROR string = "configuration_error"
	SERVICE_ERROR       string = "service_error"
	DRIVER_ERROR        string = "driver_error"
	NETWORK_ERROR       string = "network_error"
	PERMISSION_ERROR    string = "permission_error"
	CRITICAL_ERROR      string = "critical_error"
)

type SystemErrorEvent struct {
	Timestamp   time.Time              `json:"timestamp"`
	ErrorID     string                 `json:"error_id"`
	ErrorCode   string                 `json:"error_code"`
	ErrorType   string                 `json:"error_type"`
	Message     string                 `json:"message"`
	Component   string                 `json:"component"`
	Severity    string                 `json:"severity"`
	StackTrace  string                 `json:"stack_trace,omitempty"`
	Context     map[string]interface{} `json:"context"`
	UserID      string                 `json:"user_id"`
	Recoverable bool                   `json:"recoverable"`
}

func (e SystemErrorEvent) GetType() string {
	return SYSTEM_ERROR
}

func (e SystemErrorEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e SystemErrorEvent) GetData() map[string]interface{} {
	return map[string]interface{}{
		"timestamp":   e.Timestamp,
		"error_id":    e.ErrorID,
		"error_code":  e.ErrorCode,
		"error_type":  e.ErrorType,
		"message":     e.Message,
		"component":   e.Component,
		"severity":    e.Severity,
		"stack_trace": e.StackTrace,
		"context":     e.Context,
		"user_id":     e.UserID,
		"recoverable": e.Recoverable,
	}
}

type CriticalErrorEvent struct {
	Timestamp        time.Time              `json:"timestamp"`
	ErrorID          string                 `json:"error_id"`
	Message          string                 `json:"message"`
	Component        string                 `json:"component"`
	StackTrace       string                 `json:"stack_trace"`
	SystemState      map[string]interface{} `json:"system_state"`
	RecoveryAction   string                 `json:"recovery_action"`
	RequiresRestart  bool                   `json:"requires_restart"`
	DataLoss         bool                   `json:"data_loss"`
	UserID           string                 `json:"user_id"`
}

func (e CriticalErrorEvent) GetType() EventType {
	return CRITICAL_ERROR
}

func (e CriticalErrorEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e CriticalErrorEvent) GetData() map[string]interface{} {
	return map[string]interface{}{
		"timestamp":        e.Timestamp,
		"error_id":         e.ErrorID,
		"message":          e.Message,
		"component":        e.Component,
		"stack_trace":      e.StackTrace,
		"system_state":     e.SystemState,
		"recovery_action":  e.RecoveryAction,
		"requires_restart": e.RequiresRestart,
		"data_loss":        e.DataLoss,
		"user_id":          e.UserID,
	}
}

type PermissionErrorEvent struct {
	Timestamp       time.Time `json:"timestamp"`
	Operation       string    `json:"operation"`
	Resource        string    `json:"resource"`
	RequiredPermission string `json:"required_permission"`
	CurrentUser     string    `json:"current_user"`
	IsElevated      bool      `json:"is_elevated"`
	CanElevate      bool      `json:"can_elevate"`
	Suggestion      string    `json:"suggestion"`
}

func (e PermissionErrorEvent) GetType() EventType {
	return PERMISSION_ERROR
}

func (e PermissionErrorEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e PermissionErrorEvent) GetData() map[string]interface{} {
	return map[string]interface{}{
		"timestamp":           e.Timestamp,
		"operation":           e.Operation,
		"resource":            e.Resource,
		"required_permission": e.RequiredPermission,
		"current_user":        e.CurrentUser,
		"is_elevated":         e.IsElevated,
		"can_elevate":         e.CanElevate,
		"suggestion":          e.Suggestion,
	}
}
package entities

import "time"

type RegistryEntry struct {
	ID          string      `json:"id" validate:"required"`
	KeyPath     string      `json:"key_path" validate:"required"`
	ValueName   string      `json:"value_name"`
	ValueType   string      `json:"value_type" validate:"oneof=REG_SZ REG_DWORD REG_QWORD REG_BINARY REG_MULTI_SZ"`
	Value       interface{} `json:"value"`
	BackupValue interface{} `json:"backup_value"`
	IsModified  bool        `json:"is_modified"`
	IsBackedUp  bool        `json:"is_backed_up"`
	Purpose     string      `json:"purpose"`
	Category    string      `json:"category"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

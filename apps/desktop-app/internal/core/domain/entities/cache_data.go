package entities

import "time"

type CacheData struct {
	ID                    string    `json:"id" validate:"required"`
	CacheType             string    `json:"cache_type" validate:"required"`
	Location              string    `json:"location"`
	SizeBytes             int64     `json:"size_bytes"`
	FileCount             int       `json:"file_count"`
	LastCleanup           time.Time `json:"last_cleanup"`
	IsTemporary           bool      `json:"is_temporary"`
	CanBeSafelyDeleted    bool      `json:"can_be_safely_deleted"`
	CleanupScheduleEnabled bool     `json:"cleanup_schedule_enabled"`
	CleanupInterval       string    `json:"cleanup_interval"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

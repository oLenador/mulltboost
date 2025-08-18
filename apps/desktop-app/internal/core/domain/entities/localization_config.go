package entities

import "time"

type LocalizationConfig struct {
	ID                    string            `json:"id" validate:"required"`
	CurrentLanguage       string            `json:"current_language" validate:"required"`
	AvailableLanguages    []string          `json:"available_languages"`
	TranslationsPath      string            `json:"translations_path"`
	FallbackLanguage      string            `json:"fallback_language"`
	AutoDetectLanguage    bool              `json:"auto_detect_language"`
	CacheTranslations     bool              `json:"cache_translations"`
	TranslationCache      map[string]string `json:"translation_cache"`
	LastUpdated           time.Time         `json:"last_updated"`
	Version               string            `json:"version"`
	IsRTL                 bool              `json:"is_rtl"`
	DateFormat            string            `json:"date_format"`
	TimeFormat            string            `json:"time_format"`
	NumberFormat          string            `json:"number_format"`
	CurrencyFormat        string            `json:"currency_format"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

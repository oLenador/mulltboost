package i18n

import (
    "sync"
)

type Language string

const (
    English    Language = "en"
    Portuguese Language = "pt"
    Spanish    Language = "es"
)

type Translation map[string]string
type Translations map[Language]Translation

type Service struct {
    translations Translations
    mu           sync.RWMutex
}

func NewService() *Service {
    return &Service{
        translations: make(Translations),
    }
}

func (s *Service) SetTranslations(translations Translations) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.translations = translations
}

func (s *Service) Translate(key string, lang Language) string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if trans, exists := s.translations[lang]; exists {
        if value, exists := trans[key]; exists {
            return value
        }
    }
    
    if lang != English {
        if trans, exists := s.translations[English]; exists {
            if value, exists := trans[key]; exists {
                return value
            }
        }
    }
    
    return key 
}
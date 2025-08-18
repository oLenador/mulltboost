package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewSmartScreenDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_smart_screen_disable",
		NameKey:        "booster.system.smart_screen_disable.name",
		DescriptionKey: "booster.system.smart_screen_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"system", "security", "speed"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.smart_screen_disable.name":        "Отключить SmartScreen",
			"booster.system.smart_screen_disable.description": "Удаляет фильтр безопасности, ускоряя открытие файлов и программ.",
		},
		i18n.Spanish: {
			"booster.system.smart_screen_disable.name":        "Desactivar SmartScreen",
			"booster.system.smart_screen_disable.description": "Elimina el filtro de seguridad, acelerando la apertura de archivos y programas.",
		},
		i18n.Portuguese: {
			"booster.system.smart_screen_disable.name":        "Desativar SmartScreen",
			"booster.system.smart_screen_disable.description": "Remove filtro de segurança, acelerando a abertura de arquivos e programas.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.smart_screen_disable.name":        "Desativar SmartScreen",
			"booster.system.smart_screen_disable.description": "Remove filtro de segurança, acelerando a abertura de arquivos e programas.",
		},
		i18n.English: {
			"booster.system.smart_screen_disable.name":        "Disable SmartScreen",
			"booster.system.smart_screen_disable.description": "Removes the security filter, speeding up the opening of files and programs.",
		},
	}

	executor := NewSmartScreenDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewWindowsSearchDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_windows_search_disable",
		NameKey:        "booster.system.windows_search_disable.name",
		DescriptionKey: "booster.system.windows_search_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "search", "resources"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.windows_search_disable.name":        "Отключить Windows Search",
			"booster.system.windows_search_disable.description": "Отключает непрерывное индексирование, облегчая нагрузку на диск и ЦП от постоянных процессов.",
		},
		i18n.Spanish: {
			"booster.system.windows_search_disable.name":        "Desactivar Windows Search",
			"booster.system.windows_search_disable.description": "Desactiva la indexación continua, aliviando el disco y la CPU de procesos constantes.",
		},
		i18n.Portuguese: {
			"booster.system.windows_search_disable.name":        "Desativar Windows Search",
			"booster.system.windows_search_disable.description": "Desativa indexação contínua, aliviando disco e CPU de processos constantes.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.windows_search_disable.name":        "Desativar Windows Search",
			"booster.system.windows_search_disable.description": "Desativa indexação contínua, aliviando disco e CPU de processos constantes.",
		},
		i18n.English: {
			"booster.system.windows_search_disable.name":        "Disable Windows Search",
			"booster.system.windows_search_disable.description": "Disables continuous indexing, relieving the disk and CPU from constant processes.",
		},
	}

	executor := NewWindowsSearchDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewFastStartupDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_fast_startup_disable",
		NameKey:        "booster.system.fast_startup_disable.name",
		DescriptionKey: "booster.system.fast_startup_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "startup", "speed"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.fast_startup_disable.name":        "Отключить быстрый запуск",
			"booster.system.fast_startup_disable.description": "Отключает гибридную загрузку, предотвращая несоответствия при запуске системы.",
		},
		i18n.Spanish: {
			"booster.system.fast_startup_disable.name":        "Desactivar Fast Startup",
			"booster.system.fast_startup_disable.description": "Desactiva la inicialización híbrida, previniendo inconsistencias en la carga del sistema.",
		},
		i18n.Portuguese: {
			"booster.system.fast_startup_disable.name":        "Desativar Fast Startup",
			"booster.system.fast_startup_disable.description": "Desativa inicialização híbrida, prevenindo inconsistências no carregamento do sistema.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.fast_startup_disable.name":        "Desativar Fast Startup",
			"booster.system.fast_startup_disable.description": "Desativa inicialização híbrida, prevenindo inconsistências no carregamento do sistema.",
		},
		i18n.English: {
			"booster.system.fast_startup_disable.name":        "Disable Fast Startup",
			"booster.system.fast_startup_disable.description": "Disables hybrid startup, preventing inconsistencies in system loading.",
		},
	}

	executor := NewFastStartupDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
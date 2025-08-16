package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewExplorerOptimize() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_explorer_optimize",
		NameKey:        "booster.system.explorer_optimize.name",
		DescriptionKey: "booster.system.explorer_optimize.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "explorer", "speed"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.explorer_optimize.name":        "Оптимизировать Windows Explorer",
			"booster.system.explorer_optimize.description": "Настраивает параметры проводника, ускоряя навигацию и доступ к файлам.",
		},
		i18n.Spanish: {
			"booster.system.explorer_optimize.name":        "Optimizar Windows Explorer",
			"booster.system.explorer_optimize.description": "Refina las configuraciones del Explorer, acelerando la navegación y el acceso al sistema de archivos.",
		},
		i18n.Portuguese: {
			"booster.system.explorer_optimize.name":        "Otimizar Windows Explorer",
			"booster.system.explorer_optimize.description": "Refina configurações do Explorer, acelerando navegação e acesso ao sistema de arquivos.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.explorer_optimize.name":        "Otimizar Windows Explorer",
			"booster.system.explorer_optimize.description": "Refina configurações do Explorer, acelerando navegação e acesso ao sistema de arquivos.",
		},
		i18n.English: {
			"booster.system.explorer_optimize.name":        "Optimize Windows Explorer",
			"booster.system.explorer_optimize.description": "Refines Explorer settings, speeding up navigation and file system access.",
		},
	}

	executor := NewExplorerOptimizeExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
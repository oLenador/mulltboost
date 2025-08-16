package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewHibernationDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_hibernation_disable",
		NameKey:        "booster.system.hibernation_disable.name",
		DescriptionKey: "booster.system.hibernation_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "storage", "hibernation"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.hibernation_disable.name":        "Отключить спящий режим",
			"booster.system.hibernation_disable.description": "Удаляет файл гибернации, освобождая гигабайты на жестком диске.",
		},
		i18n.Spanish: {
			"booster.system.hibernation_disable.name":        "Desactivar Hibernación",
			"booster.system.hibernation_disable.description": "Elimina el archivo de hibernación, liberando gigabytes de espacio en el disco duro.",
		},
		i18n.Portuguese: {
			"booster.system.hibernation_disable.name":        "Desativar Hibernação",
			"booster.system.hibernation_disable.description": "Elimina o arquivo de hibernação, libertando gigabytes de espaço no disco rígido.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.hibernation_disable.name":        "Desativar Hibernação",
			"booster.system.hibernation_disable.description": "Elimina o arquivo de hibernação, liberando gigabytes de espaço no disco rígido.",
		},
		i18n.English: {
			"booster.system.hibernation_disable.name":        "Disable Hibernation",
			"booster.system.hibernation_disable.description": "Eliminates the hibernation file, freeing up gigabytes of space on the hard drive.",
		},
	}

	executor := NewHibernationDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
package memory

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewStandbyListClean() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "memory_standby_list_clean",
		NameKey:        "booster.memory.standby_list_clean.name",
		DescriptionKey: "booster.memory.standby_list_clean.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     false,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"memory", "standby", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.memory.standby_list_clean.name":        "Очистка списка ожидания",
			"booster.memory.standby_list_clean.description": "Удаляет ожидающие данные из памяти (экспериментально), потенциально освобождая больше ресурсов.",
		},
		i18n.Spanish: {
			"booster.memory.standby_list_clean.name":        "Limpieza de Standby List",
			"booster.memory.standby_list_clean.description": "Elimina datos en espera en la memoria (experimental), liberando potencialmente más recursos.",
		},
		i18n.Portuguese: {
			"booster.memory.standby_list_clean.name":        "Limpeza de Standby List",
			"booster.memory.standby_list_clean.description": "Remove dados em espera na memória (experimental), potencialmente libertando mais recursos.",
		},
		i18n.PortugueseBrazil: {
			"booster.memory.standby_list_clean.name":        "Limpeza de Standby List",
			"booster.memory.standby_list_clean.description": "Remove dados em espera na memória (experimental), potencialmente liberando mais recursos.",
		},
		i18n.English: {
			"booster.memory.standby_list_clean.name":        "Clean Standby List",
			"booster.memory.standby_list_clean.description": "Removes data on standby from memory (experimental), potentially freeing up more resources.",
		},
	}

	executor := NewStandbyListCleanExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
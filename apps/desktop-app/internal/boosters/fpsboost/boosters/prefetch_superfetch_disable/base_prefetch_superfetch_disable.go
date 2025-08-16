package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewPrefetchSuperfetchDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_prefetch_superfetch_disable",
		NameKey:        "booster.system.prefetch_superfetch_disable.name",
		DescriptionKey: "booster.system.prefetch_superfetch_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "memory", "prefetch"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.prefetch_superfetch_disable.name":        "Отключить Prefetch и Superfetch",
			"booster.system.prefetch_superfetch_disable.description": "Предотвращает автоматическую предварительную загрузку, освобождая ОЗУ и диск во время интенсивных задач.",
		},
		i18n.Spanish: {
			"booster.system.prefetch_superfetch_disable.name":        "Desactivar Prefetch y Superfetch",
			"booster.system.prefetch_superfetch_disable.description": "Impide la precarga automática, aliviando la RAM y el disco en tareas intensivas.",
		},
		i18n.Portuguese: {
			"booster.system.prefetch_superfetch_disable.name":        "Desativar Prefetch e Superfetch",
			"booster.system.prefetch_superfetch_disable.description": "Impede pré-carregamento automático, aliviando RAM e disco em tarefas intensivas.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.prefetch_superfetch_disable.name":        "Desativar Prefetch e Superfetch",
			"booster.system.prefetch_superfetch_disable.description": "Impede pré-carregamento automático, aliviando RAM e disco em tarefas intensivas.",
		},
		i18n.English: {
			"booster.system.prefetch_superfetch_disable.name":        "Disable Prefetch and Superfetch",
			"booster.system.prefetch_superfetch_disable.description": "Prevents automatic pre-loading, relieving RAM and disk in intensive tasks.",
		},
	}

	executor := NewPrefetchSuperfetchDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
package memory

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewNonPagedPoolOptimize() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "memory_non_paged_pool_optimize",
		NameKey:        "booster.memory.non_paged_pool_optimize.name",
		DescriptionKey: "booster.memory.non_paged_pool_optimize.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskHigh,
		Version:        "1.0.0",
		Tags:           []string{"memory", "driver", "optimization"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.memory.non_paged_pool_optimize.name":        "Оптимизировать непагинируемый пул памяти",
			"booster.memory.non_paged_pool_optimize.description": "Настраивает память, зарезервированную для драйверов, избегая узких мест в критически важных операциях.",
		},
		i18n.Spanish: {
			"booster.memory.non_paged_pool_optimize.name":        "Optimizar el Pool de Memoria No Paginação",
			"booster.memory.non_paged_pool_optimize.description": "Ajusta la memoria reservada para controladores, evitando cuellos de botella en operaciones críticas.",
		},
		i18n.Portuguese: {
			"booster.memory.non_paged_pool_optimize.name":        "Otimizar Pool de Memória Não Paginação",
			"booster.memory.non_paged_pool_optimize.description": "Ajusta memória reservada para drivers, evitando gargalos em operações críticas.",
		},
		i18n.PortugueseBrazil: {
			"booster.memory.non_paged_pool_optimize.name":        "Otimizar Pool de Memória Não Paginação",
			"booster.memory.non_paged_pool_optimize.description": "Ajusta memória reservada para drivers, evitando gargalos em operações críticas.",
		},
		i18n.English: {
			"booster.memory.non_paged_pool_optimize.name":        "Optimize Non-Paged Pool Memory",
			"booster.memory.non_paged_pool_optimize.description": "Adjusts memory reserved for drivers, preventing bottlenecks in critical operations.",
		},
	}

	executor := NewNonPagedPoolOptimizeExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
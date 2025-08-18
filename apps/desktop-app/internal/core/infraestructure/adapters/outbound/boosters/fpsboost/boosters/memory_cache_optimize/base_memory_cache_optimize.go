package memory

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewMemoryCacheOptimize() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "memory_cache_optimize",
		NameKey:        "booster.memory.cache_optimize.name",
		DescriptionKey: "booster.memory.cache_optimize.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"memory", "cache", "latency"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.memory.cache_optimize.name":        "Оптимизировать кэш ОЗУ для уменьшения задержки",
			"booster.memory.cache_optimize.description": "Улучшает временное хранение в ОЗУ, ускоряя получение часто используемых данных.",
		},
		i18n.Spanish: {
			"booster.memory.cache_optimize.name":        "Optimizar el Caché de Memoria RAM para Reducir Latencia",
			"booster.memory.cache_optimize.description": "Refina el almacenamiento temporal de la RAM, acelerando la recuperación de datos frecuentes.",
		},
		i18n.Portuguese: {
			"booster.memory.cache_optimize.name":        "Otimizar Cache de Memória RAM para Reduzir Latência",
			"booster.memory.cache_optimize.description": "Refina o armazenamento temporário da RAM, acelerando a recuperação de dados frequentes.",
		},
		i18n.PortugueseBrazil: {
			"booster.memory.cache_optimize.name":        "Otimizar Cache de Memória RAM para Reduzir Latência",
			"booster.memory.cache_optimize.description": "Refina o armazenamento temporário da RAM, acelerando a recuperação de dados frequentes.",
		},
		i18n.English: {
			"booster.memory.cache_optimize.name":        "Optimize RAM Cache to Reduce Latency",
			"booster.memory.cache_optimize.description": "Refines the temporary storage of RAM, speeding up the retrieval of frequent data.",
		},
	}

	executor := NewMemoryCacheOptimizeExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
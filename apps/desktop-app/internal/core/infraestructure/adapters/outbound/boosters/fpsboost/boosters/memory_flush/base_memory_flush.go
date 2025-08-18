package memory

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewMemoryFlush() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "memory_flush",
		NameKey:        "booster.memory.flush.name",
		DescriptionKey: "booster.memory.flush.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     false,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"memory", "flush", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.memory.flush.name":        "Очистка кэша системы",
			"booster.memory.flush.description": "Пытается очистить кэш системы, используя NtSetSystemInformation, улучшая производительность.",
		},
		i18n.Spanish: {
			"booster.memory.flush.name":        "Limpieza de Caché del Sistema",
			"booster.memory.flush.description": "Intenta limpiar el caché del sistema usando NtSetSystemInformation, mejorando el rendimiento al liberar memoria.",
		},
		i18n.Portuguese: {
			"booster.memory.flush.name":        "Limpeza de Cache do Sistema",
			"booster.memory.flush.description": "Tenta limpar o cache do sistema usando NtSetSystemInformation, melhorando o desempenho ao liberar memória.",
		},
		i18n.PortugueseBrazil: {
			"booster.memory.flush.name":        "Limpeza de Cache do Sistema",
			"booster.memory.flush.description": "Tenta limpar o cache do sistema usando NtSetSystemInformation, melhorando o desempenho ao liberar memória.",
		},
		i18n.English: {
			"booster.memory.flush.name":        "System Cache Cleanup",
			"booster.memory.flush.description": "Attempts to clean the system cache using NtSetSystemInformation, improving performance by freeing up memory.",
		},
	}

	executor := NewMemoryFlushExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
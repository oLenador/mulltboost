package memory

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewMemoryAllocation() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "memory_allocation",
		NameKey:        "booster.memory.allocation.name",
		DescriptionKey: "booster.memory.allocation.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskHigh,
		Version:        "1.0.0",
		Tags:           []string{"memory", "ram", "optimization"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.memory.allocation.name":        "Оптимизировать выделение памяти",
			"booster.memory.allocation.description": "Резервирует ОЗУ для приоритетных процессов, минимизируя фрагментацию и общую задержку.",
		},
		i18n.Spanish: {
			"booster.memory.allocation.name":        "Optimizar la Asignación de Memoria",
			"booster.memory.allocation.description": "Reserva RAM para procesos prioritarios, minimizando la fragmentación y la latencia general.",
		},
		i18n.Portuguese: {
			"booster.memory.allocation.name":        "Otimizar Alocação de Memória",
			"booster.memory.allocation.description": "Reserva RAM para processos prioritários, minimizando fragmentação e latência geral.",
		},
		i18n.PortugueseBrazil: {
			"booster.memory.allocation.name":        "Otimizar Alocação de Memória",
			"booster.memory.allocation.description": "Reserva RAM para processos prioritários, minimizando fragmentação e latência geral.",
		},
		i18n.English: {
			"booster.memory.allocation.name":        "Optimize Memory Allocation",
			"booster.memory.allocation.description": "Reserves RAM for priority processes, minimizing fragmentation and overall latency.",
		},
	}

	executor := NewMemoryAllocationExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
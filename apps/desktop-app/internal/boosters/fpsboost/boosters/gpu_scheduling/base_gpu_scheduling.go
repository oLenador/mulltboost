package gpu

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewGPUScheduling() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "gpu_scheduling",
		NameKey:        "booster.gpu.scheduling.name",
		DescriptionKey: "booster.gpu.scheduling.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"gpu", "scheduling", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.gpu.scheduling.name":        "Оптимизировать планирование GPU",
			"booster.gpu.scheduling.description": "Реорганизует графические очереди, ускоряя обработку рендеринга в реальном времени.",
		},
		i18n.Spanish: {
			"booster.gpu.scheduling.name":        "Optimizar el Agendamiento de GPU",
			"booster.gpu.scheduling.description": "Reorganiza las colas gráficas, acelerando el procesamiento de renderización en tiempo real.",
		},
		i18n.Portuguese: {
			"booster.gpu.scheduling.name":        "Otimizar Agendamento de GPU",
			"booster.gpu.scheduling.description": "Reorganiza filas gráficas, acelerando processamento de renderização em tempo real.",
		},
		i18n.PortugueseBrazil: {
			"booster.gpu.scheduling.name":        "Otimizar Agendamento de GPU",
			"booster.gpu.scheduling.description": "Reorganiza filas gráficas, acelerando processamento de renderização em tempo real.",
		},
		i18n.English: {
			"booster.gpu.scheduling.name":        "Optimize GPU Scheduling",
			"booster.gpu.scheduling.description": "Reorganizes graphical queues, speeding up real-time rendering processing.",
		},
	}

	executor := NewGPUSchedulingExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
package gpu

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewGPUIRQOptimize() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "gpu_irq_optimize",
		NameKey:        "booster.gpu.irq_optimize.name",
		DescriptionKey: "booster.gpu.irq_optimize.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskHigh,
		Version:        "1.0.0",
		Tags:           []string{"gpu", "irq", "latency"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.gpu.irq_optimize.name":        "Оптимизировать IRQ GPU",
			"booster.gpu.irq_optimize.description": "Повышает приоритет графических прерываний, ускоряя отклик GPU.",
		},
		i18n.Spanish: {
			"booster.gpu.irq_optimize.name":        "Optimizar IRQ de la GPU",
			"booster.gpu.irq_optimize.description": "Eleva la prioridad de interrupciones gráficas, agilizando las respuestas de la GPU.",
		},
		i18n.Portuguese: {
			"booster.gpu.irq_optimize.name":        "Otimizar IRQ da GPU",
			"booster.gpu.irq_optimize.description": "Eleva a prioridade de interrupções gráficas, agilizando as respostas da GPU.",
		},
		i18n.PortugueseBrazil: {
			"booster.gpu.irq_optimize.name":        "Otimizar IRQ da GPU",
			"booster.gpu.irq_optimize.description": "Eleva a prioridade de interrupções gráficas, agilizando as respostas da GPU.",
		},
		i18n.English: {
			"booster.gpu.irq_optimize.name":        "Optimize GPU IRQ",
			"booster.gpu.irq_optimize.description": "Elevates the priority of graphical interruptions, speeding up GPU responses.",
		},
	}

	executor := NewGPUIRQOptimizeExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
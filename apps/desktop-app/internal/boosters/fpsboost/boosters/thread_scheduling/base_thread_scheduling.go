package cpu

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewThreadScheduling() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "cpu_thread_scheduling",
		NameKey:        "booster.cpu.thread_scheduling.name",
		DescriptionKey: "booster.cpu.thread_scheduling.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"cpu", "threads", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.cpu.thread_scheduling.name":        "Оптимизировать планирование потоков",
			"booster.cpu.thread_scheduling.description": "Улучшает распределение потоков, максимизируя эффективное использование ядер ЦП.",
		},
		i18n.Spanish: {
			"booster.cpu.thread_scheduling.name":        "Optimizar el Agendamiento de Threads",
			"booster.cpu.thread_scheduling.description": "Refina la distribución de threads, maximizando el uso eficiente de los núcleos de la CPU.",
		},
		i18n.Portuguese: {
			"booster.cpu.thread_scheduling.name":        "Otimizar Agendamento de Threads",
			"booster.cpu.thread_scheduling.description": "Refina a distribuição de threads, maximizando o uso eficiente dos núcleos da CPU.",
		},
		i18n.PortugueseBrazil: {
			"booster.cpu.thread_scheduling.name":        "Otimizar Agendamento de Threads",
			"booster.cpu.thread_scheduling.description": "Refina a distribuição de threads, maximizando o uso eficiente dos núcleos da CPU.",
		},
		i18n.English: {
			"booster.cpu.thread_scheduling.name":        "Optimize Thread Scheduling",
			"booster.cpu.thread_scheduling.description": "Refines thread distribution, maximizing the efficient use of CPU cores.",
		},
	}

	executor := NewThreadSchedulingExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
package cpu

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewCPUAffinityBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "cpu_affinity",
		NameKey:        "booster.cpu.affinity.name",
		DescriptionKey: "booster.cpu.affinity.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"cpu", "performance", "optimization"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.cpu.affinity.name":        "Настроить аффинность ЦП",
			"booster.cpu.affinity.description": "Направляет задачи на конкретные ядра, балансируя нагрузку и повышая эффективность.",
		},
		i18n.Spanish: {
			"booster.cpu.affinity.name":        "Ajustar Afinidad de la CPU",
			"booster.cpu.affinity.description": "Dirige las tareas a núcleos específicos, equilibrando la carga y aumentando la eficiencia.",
		},
		i18n.Portuguese: {
			"booster.cpu.affinity.name":        "Ajustar Afinidade da CPU",
			"booster.cpu.affinity.description": "Direciona tarefas a núcleos específicos, balanceando carga e aumentando eficiência.",
		},
		i18n.PortugueseBrazil: {
			"booster.cpu.affinity.name":        "Ajustar Afinidade da CPU",
			"booster.cpu.affinity.description": "Direciona tarefas a núcleos específicos, balanceando carga e aumentando eficiência.",
		},
		i18n.English: {
			"booster.cpu.affinity.name":        "Adjust CPU Affinity",
			"booster.cpu.affinity.description": "Directs tasks to specific cores, balancing the load and increasing efficiency.",
		},
	}

	executor := NewCPUAffinityExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
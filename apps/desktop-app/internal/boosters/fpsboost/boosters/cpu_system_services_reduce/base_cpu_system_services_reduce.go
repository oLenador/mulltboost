package cpu

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewCPUSystemServicesReduce() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "cpu_system_services_reduce",
		NameKey:        "booster.cpu.system_services_reduce.name",
		DescriptionKey: "booster.cpu.system_services_reduce.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"cpu", "services", "optimization"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.cpu.system_services_reduce.name":        "Уменьшить использование ЦП системными службами",
			"booster.cpu.system_services_reduce.description": "Понижает приоритет собственных служб, направляя больше мощности ЦП на другие задачи.",
		},
		i18n.Spanish: {
			"booster.cpu.system_services_reduce.name":        "Reducir el Uso de CPU por Servicios del Sistema",
			"booster.cpu.system_services_reduce.description": "Baja la prioridad de los servicios nativos, dirigiendo más poder de la CPU.",
		},
		i18n.Portuguese: {
			"booster.cpu.system_services_reduce.name":        "Reduzir Uso de CPU por Serviços de Sistema",
			"booster.cpu.system_services_reduce.description": "Rebaixa a prioridade de serviços nativos, direcionando mais poder da CPU.",
		},
		i18n.PortugueseBrazil: {
			"booster.cpu.system_services_reduce.name":        "Reduzir Uso de CPU por Serviços de Sistema",
			"booster.cpu.system_services_reduce.description": "Rebaixa a prioridade de serviços nativos, direcionando mais poder da CPU.",
		},
		i18n.English: {
			"booster.cpu.system_services_reduce.name":        "Reduce CPU Usage by System Services",
			"booster.cpu.system_services_reduce.description": "Lowers the priority of native services, directing more CPU power elsewhere.",
		},
	}

	executor := NewCPUSystemServicesReduceExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
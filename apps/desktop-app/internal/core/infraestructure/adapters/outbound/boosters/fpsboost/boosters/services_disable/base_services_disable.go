package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewServicesDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_services_disable",
		NameKey:        "booster.system.services_disable.name",
		DescriptionKey: "booster.system.services_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskHigh,
		Version:        "1.0.0",
		Tags:           []string{"system", "services", "optimization"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.services_disable.name":        "Отключить ненужные службы",
			"booster.system.services_disable.description": "Останавливает вторичные процессы Windows, освобождая дополнительную ОЗУ и ЦП.",
		},
		i18n.Spanish: {
			"booster.system.services_disable.name":        "Desactivar Servicios Innecesarios",
			"booster.system.services_disable.description": "Interrumpe procesos secundarios de Windows, liberando RAM y CPU adicionales.",
		},
		i18n.Portuguese: {
			"booster.system.services_disable.name":        "Desativar Serviços Desnecessários",
			"booster.system.services_disable.description": "Interrompe processos secundários do Windows, libertando RAM e CPU extras.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.services_disable.name":        "Desativar Serviços Desnecessários",
			"booster.system.services_disable.description": "Interrompe processos secundários do Windows, liberando RAM e CPU extras.",
		},
		i18n.English: {
			"booster.system.services_disable.name":        "Disable Unnecessary Services",
			"booster.system.services_disable.description": "Stops secondary Windows processes, freeing up extra RAM and CPU.",
		},
	}

	executor := NewServicesDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
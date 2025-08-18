package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewXboxServicesDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_xbox_services_disable",
		NameKey:        "booster.system.xbox_services_disable.name",
		DescriptionKey: "booster.system.xbox_services_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "xbox", "gaming"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.xbox_services_disable.name":        "Отключить службы Xbox",
			"booster.system.xbox_services_disable.description": "Останавливает собственные службы консоли, уменьшая нагрузку на память и ЦП.",
		},
		i18n.Spanish: {
			"booster.system.xbox_services_disable.name":        "Desactivar Servicios de Xbox",
			"booster.system.xbox_services_disable.description": "Detiene los servicios nativos de la consola, reduciendo la carga en memoria y CPU.",
		},
		i18n.Portuguese: {
			"booster.system.xbox_services_disable.name":        "Desativar Xbox Services",
			"booster.system.xbox_services_disable.description": "Para serviços de consola nativos, reduzindo carga em memória e CPU.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.xbox_services_disable.name":        "Desativar Xbox Services",
			"booster.system.xbox_services_disable.description": "Para serviços de console nativos, reduzindo carga em memória e CPU.",
		},
		i18n.English: {
			"booster.system.xbox_services_disable.name":        "Disable Xbox Services",
			"booster.system.xbox_services_disable.description": "Stops native console services, reducing the load on memory and CPU.",
		},
	}

	executor := NewXboxServicesDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}

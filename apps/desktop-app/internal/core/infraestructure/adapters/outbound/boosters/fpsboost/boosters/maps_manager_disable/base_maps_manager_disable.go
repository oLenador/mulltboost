package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewMapsManagerDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_maps_manager_disable",
		NameKey:        "booster.system.maps_manager_disable.name",
		DescriptionKey: "booster.system.maps_manager_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "maps", "resources"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.maps_manager_disable.name":        "Отключить Download Maps Manager",
			"booster.system.maps_manager_disable.description": "Останавливает службу автономных карт, снижая использование диска и сети в фоновом режиме.",
		},
		i18n.Spanish: {
			"booster.system.maps_manager_disable.name":        "Desactivar Download Maps Manager",
			"booster.system.maps_manager_disable.description": "Detiene el servicio de mapas sin conexión, reduciendo el uso de disco y red en segundo plano.",
		},
		i18n.Portuguese: {
			"booster.system.maps_manager_disable.name":        "Desativar Download Maps Manager",
			"booster.system.maps_manager_disable.description": "Para serviço de mapas offline, reduzindo o uso de disco e rede em segundo plano.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.maps_manager_disable.name":        "Desativar Download Maps Manager",
			"booster.system.maps_manager_disable.description": "Para serviço de mapas offline, reduzindo o uso de disco e rede em segundo plano.",
		},
		i18n.English: {
			"booster.system.maps_manager_disable.name":        "Disable Download Maps Manager",
			"booster.system.maps_manager_disable.description": "Stops the offline maps service, reducing disk and network usage in the background.",
		},
	}

	executor := NewMapsManagerDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
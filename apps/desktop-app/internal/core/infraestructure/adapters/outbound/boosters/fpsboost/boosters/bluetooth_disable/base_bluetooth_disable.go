package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewBluetoothDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_bluetooth_disable",
		NameKey:        "booster.system.bluetooth_disable.name",
		DescriptionKey: "booster.system.bluetooth_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "bluetooth", "resources"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.bluetooth_disable.name":        "Отключить Bluetooth",
			"booster.system.bluetooth_disable.description": "Отключает службу Bluetooth, освобождая ресурсы и уменьшая системные прерывания.",
		},
		i18n.Spanish: {
			"booster.system.bluetooth_disable.name":        "Desactivar Bluetooth",
			"booster.system.bluetooth_disable.description": "Desactiva el servicio de Bluetooth, liberando recursos y reduciendo las interrupciones del sistema.",
		},
		i18n.Portuguese: {
			"booster.system.bluetooth_disable.name":        "Desativar Bluetooth",
			"booster.system.bluetooth_disable.description": "Desliga o serviço de Bluetooth, libertando recursos e reduzindo interrupções do sistema.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.bluetooth_disable.name":        "Desativar Bluetooth",
			"booster.system.bluetooth_disable.description": "Desliga o serviço de Bluetooth, liberando recursos e reduzindo interrupções do sistema.",
		},
		i18n.English: {
			"booster.system.bluetooth_disable.name":        "Disable Bluetooth",
			"booster.system.bluetooth_disable.description": "Turns off the Bluetooth service, freeing up resources and reducing system interruptions.",
		},
	}

	executor := NewBluetoothDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
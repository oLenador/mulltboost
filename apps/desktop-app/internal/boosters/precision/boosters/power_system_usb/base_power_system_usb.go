package precision

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewPowerSystemUSBBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "precision_power_system_usb",
		NameKey:        "booster.precision.power_system_usb.name",
		DescriptionKey: "booster.precision.power_system_usb.description",
		Category:       entities.CategoryPrecision,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"power", "system", "usb", "latency"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.precision.power_system_usb.name":        "Оптимизация питания, системы и USB",
			"booster.precision.power_system_usb.description": "Приоритезирует производительность по питанию и USB, отключая таймеры и оптимизируя задержку портов.",
		},
		i18n.Spanish: {
			"booster.precision.power_system_usb.name":        "Optimizaciones de Energía, Sistema y USB",
			"booster.precision.power_system_usb.description": "Prioriza el rendimiento en energía y USB, desactivando temporizadores, mejorando la accesibilidad y optimizando la latencia de los puertos.",
		},
		i18n.Portuguese: {
			"booster.precision.power_system_usb.name":        "Otimizações de Energia, Sistema e USB",
			"booster.precision.power_system_usb.description": "Prioriza o desempenho em energia e USB, desativando temporizadores, melhorando acessibilidade e otimizando a latência de portas.",
		},
		i18n.PortugueseBrazil: {
			"booster.precision.power_system_usb.name":        "Otimizações de Energia, Sistema e USB",
			"booster.precision.power_system_usb.description": "Prioriza o desempenho em energia e USB, desativando temporizadores, melhorando acessibilidade e otimizando a latência de portas.",
		},
		i18n.English: {
			"booster.precision.power_system_usb.name":        "Power, System, and USB Optimizations",
			"booster.precision.power_system_usb.description": "Prioritizes performance in power and USB, disabling timers, improving accessibility, and optimizing port latency.",
		},
	}

	executor := NewPowerSystemUSBExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
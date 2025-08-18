package precision

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewControllerPrecisionBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "precision_controller",
		NameKey:        "booster.precision.controller.name",
		DescriptionKey: "booster.precision.controller.description",
		Category:       entities.CategoryPrecision,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"controller", "precision", "usb"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.precision.controller.name":        "Оптимизация точности контроллеров",
			"booster.precision.controller.description": "Приоритезирует производительность по питанию и USB, отключая таймеры и оптимизируя задержку портов.",
		},
		i18n.Spanish: {
			"booster.precision.controller.name":        "Optimizaciones de Precisión de Controles",
			"booster.precision.controller.description": "Prioriza el rendimiento en energía y USB, desactivando temporizadores, mejorando la accesibilidad y optimizando la latencia de los puertos.",
		},
		i18n.Portuguese: {
			"booster.precision.controller.name":        "Otimizações de Precisão de Controlos",
			"booster.precision.controller.description": "Prioriza o desempenho em energia e USB, desativando temporizadores, melhorando acessibilidade e otimizando a latência de portas.",
		},
		i18n.PortugueseBrazil: {
			"booster.precision.controller.name":        "Otimizações de Precisão de Controles",
			"booster.precision.controller.description": "Prioriza o desempenho em energia e USB, desativando temporizadores, melhorando acessibilidade e otimizando a latência de portas.",
		},
		i18n.English: {
			"booster.precision.controller.name":        "Controller Precision Optimizations",
			"booster.precision.controller.description": "Prioritizes performance in power and USB, disabling timers, improving accessibility, and optimizing port latency.",
		},
	}

	executor := NewControllerPrecisionExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
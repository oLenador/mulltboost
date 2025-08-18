package precision

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewDisplayPrecisionBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "precision_display",
		NameKey:        "booster.precision.display.name",
		DescriptionKey: "booster.precision.display.description",
		Category:       entities.CategoryPrecision,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"display", "graphics", "latency"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.precision.display.name":        "Оптимизация точности графики и монитора",
			"booster.precision.display.description": "Снижает задержку графики и монитора, настраивая переходы питания для повышения визуальной плавности.",
		},
		i18n.Spanish: {
			"booster.precision.display.name":        "Optimizaciones de Gráficos y Precisión del Monitor",
			"booster.precision.display.description": "Reduce la latencia de gráficos y monitores, ajustando las transiciones de energía para una mayor fluidez visual.",
		},
		i18n.Portuguese: {
			"booster.precision.display.name":        "Otimizações de Gráficos e Precisão do Monitor",
			"booster.precision.display.description": "Reduz a latência de gráficos e monitores, ajustando transições de energia para maior fluidez visual.",
		},
		i18n.PortugueseBrazil: {
			"booster.precision.display.name":        "Otimizações de Gráficos e Precisão do Monitor",
			"booster.precision.display.description": "Reduz a latência de gráficos e monitores, ajustando transições de energia para maior fluidez visual.",
		},
		i18n.English: {
			"booster.precision.display.name":        "Graphics and Monitor Precision Optimizations",
			"booster.precision.display.description": "Reduces graphics and monitor latency by adjusting power transitions for greater visual fluidity.",
		},
	}

	executor := NewDisplayPrecisionExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
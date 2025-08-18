package precision

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewMousePrecisionBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "precision_mouse",
		NameKey:        "booster.precision.mouse.name",
		DescriptionKey: "booster.precision.mouse.description",
		Category:       entities.CategoryPrecision,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"mouse", "precision", "input"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.precision.mouse.name":        "Оптимизация точности мыши",
			"booster.precision.mouse.description": "Улучшает точность и отзывчивость мыши, настраивая чувствительность, ускорение и приоритет потоков.",
		},
		i18n.Spanish: {
			"booster.precision.mouse.name":        "Optimizaciones de Precisión del Ratón",
			"booster.precision.mouse.description": "Mejora la precisión y la capacidad de respuesta del ratón, ajustando la sensibilidad, la aceleración, el magnetismo, la cola de datos y priorizando los hilos.",
		},
		i18n.Portuguese: {
			"booster.precision.mouse.name":        "Otimizações de Precisão do Rato",
			"booster.precision.mouse.description": "Melhora a precisão e a capacidade de resposta do rato, ajustando sensibilidade, aceleração, magnetismo, fila de dados e priorizando threads.",
		},
		i18n.PortugueseBrazil: {
			"booster.precision.mouse.name":        "Otimizações de Precisão do Mouse",
			"booster.precision.mouse.description": "Melhora a precisão e a capacidade de resposta do mouse, ajustando sensibilidade, aceleração, magnetismo, fila de dados e priorizando threads.",
		},
		i18n.English: {
			"booster.precision.mouse.name":        "Mouse Precision Optimizations",
			"booster.precision.mouse.description": "Improves mouse precision and responsiveness by adjusting sensitivity, acceleration, magnetism, data queue, and prioritizing threads.",
		},
	}

	executor := NewMousePrecisionExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
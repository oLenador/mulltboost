package precision

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewKeyboardPrecisionBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "precision_keyboard",
		NameKey:        "booster.precision.keyboard.name",
		DescriptionKey: "booster.precision.keyboard.description",
		Category:       entities.CategoryPrecision,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"keyboard", "precision", "input"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.precision.keyboard.name":        "Оптимизация точности клавиатуры",
			"booster.precision.keyboard.description": "Улучшает доступность и отзывчивость клавиатуры, устраняя задержки и оптимизируя приоритет.",
		},
		i18n.Spanish: {
			"booster.precision.keyboard.name":        "Optimizaciones de Precisión del Teclado",
			"booster.precision.keyboard.description": "Mejora la accesibilidad y la capacidad de respuesta del teclado, eliminando retrasos, ajustando la cola de datos y optimizando la prioridad.",
		},
		i18n.Portuguese: {
			"booster.precision.keyboard.name":        "Otimizações de Precisão do Teclado",
			"booster.precision.keyboard.description": "Melhora a acessibilidade e a capacidade de resposta do teclado, removendo atrasos, ajustando fila de dados e otimizando prioridade.",
		},
		i18n.PortugueseBrazil: {
			"booster.precision.keyboard.name":        "Otimizações de Precisão do Teclado",
			"booster.precision.keyboard.description": "Melhora a acessibilidade e a capacidade de resposta do teclado, removendo atrasos, ajustando fila de dados e otimizando prioridade.",
		},
		i18n.English: {
			"booster.precision.keyboard.name":        "Keyboard Precision Optimizations",
			"booster.precision.keyboard.description": "Improves keyboard accessibility and responsiveness by removing delays, adjusting data queue, and optimizing priority.",
		},
	}

	executor := NewKeyboardPrecisionExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
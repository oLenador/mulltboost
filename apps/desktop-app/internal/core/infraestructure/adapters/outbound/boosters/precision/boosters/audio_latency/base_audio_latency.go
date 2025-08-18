package precision

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewAudioLatencyBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "precision_audio",
		NameKey:        "booster.precision.audio.name",
		DescriptionKey: "booster.precision.audio.description",
		Category:       entities.CategoryPrecision,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"audio", "latency", "sound"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.precision.audio.name":        "Оптимизация для уменьшения задержки аудио",
			"booster.precision.audio.description": "Уменьшает задержку аудио, настраивая параметры для большей плавности звука.",
		},
		i18n.Spanish: {
			"booster.precision.audio.name":        "Optimizaciones para Disminuir la Latencia de Audio",
			"booster.precision.audio.description": "Reduce la latencia de audio, ajustando las transiciones de energía para una mayor fluidez de sonido.",
		},
		i18n.Portuguese: {
			"booster.precision.audio.name":        "Otimizações para Diminuir Latência de Áudio",
			"booster.precision.audio.description": "Reduz a latência de áudio, ajustando transições de energia para maior fluidez de som.",
		},
		i18n.PortugueseBrazil: {
			"booster.precision.audio.name":        "Otimizações para Diminuir Latência de Áudio",
			"booster.precision.audio.description": "Reduz a latência de áudio, ajustando transições de energia para maior fluidez de som.",
		},
		i18n.English: {
			"booster.precision.audio.name":        "Audio Latency Reduction Optimizations",
			"booster.precision.audio.description": "Reduces audio latency by adjusting power transitions for greater sound fluidity.",
		},
	}

	executor := NewAudioLatencyExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
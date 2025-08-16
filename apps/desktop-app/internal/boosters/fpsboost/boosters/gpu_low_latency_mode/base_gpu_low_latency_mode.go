package gpu

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewGPULowLatencyMode() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "gpu_low_latency_mode",
		NameKey:        "booster.gpu.low_latency_mode.name",
		DescriptionKey: "booster.gpu.low_latency_mode.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"gpu", "latency", "graphics"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.gpu.low_latency_mode.name":        "Включить режим низкой задержки на GPU",
			"booster.gpu.low_latency_mode.description": "Настраивает драйвер GPU для минимизации задержек при интенсивной графической обработке.",
		},
		i18n.Spanish: {
			"booster.gpu.low_latency_mode.name":        "Habilitar Modo de Baja Latencia en la GPU",
			"booster.gpu.low_latency_mode.description": "Configura el controlador de la GPU para minimizar los retrasos en el procesamiento gráfico intenso.",
		},
		i18n.Portuguese: {
			"booster.gpu.low_latency_mode.name":        "Habilitar Modo de Baixa Latência na GPU",
			"booster.gpu.low_latency_mode.description": "Configura o driver da GPU para minimizar atrasos em processamento gráfico intenso.",
		},
		i18n.PortugueseBrazil: {
			"booster.gpu.low_latency_mode.name":        "Habilitar Modo de Baixa Latência na GPU",
			"booster.gpu.low_latency_mode.description": "Configura o driver da GPU para minimizar atrasos em processamento gráfico intenso.",
		},
		i18n.English: {
			"booster.gpu.low_latency_mode.name":        "Enable Low Latency Mode on GPU",
			"booster.gpu.low_latency_mode.description": "Configures the GPU driver to minimize delays in intense graphics processing.",
		},
	}

	executor := NewGPULowLatencyModeExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
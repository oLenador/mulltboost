package gpu

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewGPUPowerPolicy() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "gpu_power_policy",
		NameKey:        "booster.gpu.power_policy.name",
		DescriptionKey: "booster.gpu.power_policy.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"gpu", "power", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.gpu.power_policy.name":        "Настроить политику питания GPU",
			"booster.gpu.power_policy.description": "Поддерживает GPU в режиме максимальной производительности, исключая падения из-за экономии энергии.",
		},
		i18n.Spanish: {
			"booster.gpu.power_policy.name":        "Ajustar Política de Energía de la GPU",
			"booster.gpu.power_policy.description": "Mantiene la GPU en su máximo rendimiento, eliminando caídas por ahorro de energía.",
		},
		i18n.Portuguese: {
			"booster.gpu.power_policy.name":        "Ajustar Política de Energia da GPU",
			"booster.gpu.power_policy.description": "Mantém GPU em desempenho máximo, eliminando quedas por economia energética.",
		},
		i18n.PortugueseBrazil: {
			"booster.gpu.power_policy.name":        "Ajustar Política de Energia da GPU",
			"booster.gpu.power_policy.description": "Mantém GPU em desempenho máximo, eliminando quedas por economia energética.",
		},
		i18n.English: {
			"booster.gpu.power_policy.name":        "Adjust GPU Power Policy",
			"booster.gpu.power_policy.description": "Keeps the GPU at maximum performance, eliminating drops due to power saving.",
		},
	}

	executor := NewGPUPowerPolicyExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
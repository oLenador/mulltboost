package power

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewMultBoostPower() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "power_mult_boost_power",
		NameKey:        "booster.power.mult_boost_power.name",
		DescriptionKey: "booster.power.mult_boost_power.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"power", "performance", "energy"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.power.mult_boost_power.name":        "Mult Boost Power",
			"booster.power.mult_boost_power.description": "Автоматически настраивает систему для максимальной производительности или энергоэффективности.",
		},
		i18n.Spanish: {
			"booster.power.mult_boost_power.name":        "Mult Boost Power",
			"booster.power.mult_boost_power.description": "Configura automáticamente el sistema para operar en modos que favorecen el rendimiento máximo o la eficiencia energética.",
		},
		i18n.Portuguese: {
			"booster.power.mult_boost_power.name":        "Mult Boost Power",
			"booster.power.mult_boost_power.description": "Configura automaticamente o sistema para operar em modos que favorecem o desempenho máximo ou a eficiência energética do computador.",
		},
		i18n.PortugueseBrazil: {
			"booster.power.mult_boost_power.name":        "Mult Boost Power",
			"booster.power.mult_boost_power.description": "Configura automaticamente o sistema para operar em modos que favorecem o desempenho máximo ou a eficiência energética do computador.",
		},
		i18n.English: {
			"booster.power.mult_boost_power.name":        "Mult Boost Power",
			"booster.power.mult_boost_power.description": "Automatically configures the system to operate in modes that favor maximum performance or energy efficiency.",
		},
	}

	executor := NewMultBoostPowerExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
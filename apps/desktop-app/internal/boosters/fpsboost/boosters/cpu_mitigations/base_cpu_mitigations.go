package cpu

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewCPUMitigationsBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "cpu_mitigations",
		NameKey:        "booster.cpu.mitigations.name",
		DescriptionKey: "booster.cpu.mitigations.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskHigh,
		Version:        "1.0.0",
		Tags:           []string{"cpu", "security", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.cpu.mitigations.name":        "Оптимизировать митигации ЦП (Spectre/Meltdown)",
			"booster.cpu.mitigations.description": "Удаляет барьеры безопасности ЦП, повышая производительность при высоких нагрузках.",
		},
		i18n.Spanish: {
			"booster.cpu.mitigations.name":        "Optimizar Mitigaciones de CPU (Spectre/Meltdown)",
			"booster.cpu.mitigations.description": "Elimina las barreras de seguridad de la CPU, aumentando el rendimiento en cargas pesadas.",
		},
		i18n.Portuguese: {
			"booster.cpu.mitigations.name":        "Otimizar Mitigações de CPU (Spectre/Meltdown)",
			"booster.cpu.mitigations.description": "Remove barreiras de segurança da CPU, elevando desempenho em cargas pesadas.",
		},
		i18n.PortugueseBrazil: {
			"booster.cpu.mitigations.name":        "Otimizar Mitigações de CPU (Spectre/Meltdown)",
			"booster.cpu.mitigations.description": "Remove barreiras de segurança da CPU, elevando desempenho em cargas pesadas.",
		},
		i18n.English: {
			"booster.cpu.mitigations.name":        "Optimize CPU Mitigations (Spectre/Meltdown)",
			"booster.cpu.mitigations.description": "Removes CPU security barriers, boosting performance under heavy loads.",
		},
	}

	executor := NewCPUMitigationsExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
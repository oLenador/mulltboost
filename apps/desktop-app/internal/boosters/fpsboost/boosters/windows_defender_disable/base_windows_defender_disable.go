package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewWindowsDefenderDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_windows_defender_disable",
		NameKey:        "booster.system.windows_defender_disable.name",
		DescriptionKey: "booster.system.windows_defender_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskHigh,
		Version:        "1.0.0",
		Tags:           []string{"system", "defender", "security"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.windows_defender_disable.name":        "Отключить Windows Defender в реальном времени",
			"booster.system.windows_defender_disable.description": "Останавливает постоянное сканирование, уменьшая потребление ЦП и доступ к диску.",
		},
		i18n.Spanish: {
			"booster.system.windows_defender_disable.name":        "Desactivar Windows Defender en Tiempo Real",
			"booster.system.windows_defender_disable.description": "Detiene los escaneos constantes, disminuyendo el consumo de CPU y el acceso al disco.",
		},
		i18n.Portuguese: {
			"booster.system.windows_defender_disable.name":        "Desativar Windows Defender em Tempo Real",
			"booster.system.windows_defender_disable.description": "Para varreduras constantes, diminuindo consumo de CPU e acesso ao disco.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.windows_defender_disable.name":        "Desativar Windows Defender em Tempo Real",
			"booster.system.windows_defender_disable.description": "Para varreduras constantes, diminuindo consumo de CPU e acesso ao disco.",
		},
		i18n.English: {
			"booster.system.windows_defender_disable.name":        "Disable Windows Defender Real-Time",
			"booster.system.windows_defender_disable.description": "Stops constant scans, reducing CPU consumption and disk access.",
		},
	}

	executor := NewWindowsDefenderDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
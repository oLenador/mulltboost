package cpu

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewTimerResolution() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "cpu_timer_resolution",
		NameKey:        "booster.cpu.timer_resolution.name",
		DescriptionKey: "booster.cpu.timer_resolution.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"cpu", "timer", "latency"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.cpu.timer_resolution.name":        "Timer Resolution для улучшения отклика ЦП",
			"booster.cpu.timer_resolution.description": "Повышает точность таймера ЦП, уменьшая задержки в критических процессах.",
		},
		i18n.Spanish: {
			"booster.cpu.timer_resolution.name":        "Timer Resolution para Mejorar la Respuesta de la CPU",
			"booster.cpu.timer_resolution.description": "Aumenta la precisión del temporizador de la CPU, reduciendo los retrasos en procesos críticos.",
		},
		i18n.Portuguese: {
			"booster.cpu.timer_resolution.name":        "Timer Resolution para Melhorar Resposta da CPU",
			"booster.cpu.timer_resolution.description": "Aumenta a precisão do temporizador da CPU, reduzindo atrasos em processos críticos.",
		},
		i18n.PortugueseBrazil: {
			"booster.cpu.timer_resolution.name":        "Timer Resolution para Melhorar Resposta da CPU",
			"booster.cpu.timer_resolution.description": "Aumenta a precisão do temporizador da CPU, reduzindo atrasos em processos críticos.",
		},
		i18n.English: {
			"booster.cpu.timer_resolution.name":        "Timer Resolution to Improve CPU Response",
			"booster.cpu.timer_resolution.description": "Increases the precision of the CPU timer, reducing delays in critical processes.",
		},
	}

	executor := NewTimerResolutionExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}

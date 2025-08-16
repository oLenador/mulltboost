package cpu

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewCPUPriorityGame() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "cpu_game_priority",
		NameKey:        "booster.cpu.game_priority.name",
		DescriptionKey: "booster.cpu.game_priority.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"cpu", "games", "priority", "fps"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.cpu.game_priority.name":        "Повысить приоритет ЦП для игр",
			"booster.cpu.game_priority.description": "Настраивает приоритет обработки графики и оптимизирует использование ЦП для улучшения частоты кадров в играх.",
		},
		i18n.Spanish: {
			"booster.cpu.game_priority.name":        "Aumentar Prioridad de la CPU para Juegos",
			"booster.cpu.game_priority.description": "Ajusta la prioridad de procesamiento gráfico y optimiza el uso de la CPU para mejorar la tasa de cuadros en juegos.",
		},
		i18n.Portuguese: {
			"booster.cpu.game_priority.name":        "Aumentar Prioridade da CPU para Jogos",
			"booster.cpu.game_priority.description": "Ajusta a prioridade de processamento gráfico e otimiza o uso da CPU para melhorar a taxa de quadros em jogos.",
		},
		i18n.PortugueseBrazil: {
			"booster.cpu.game_priority.name":        "Aumentar Prioridade da CPU para Jogos",
			"booster.cpu.game_priority.description": "Ajusta a prioridade de processamento gráfico e otimiza o uso da CPU para melhorar a taxa de quadros em jogos.",
		},
		i18n.English: {
			"booster.cpu.game_priority.name":        "Increase CPU Priority for Games",
			"booster.cpu.game_priority.description": "Adjusts graphics processing priority and optimizes CPU usage to improve frame rates in games.",
		},
	}

	executor := NewCPUPriorityGameExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
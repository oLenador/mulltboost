package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewVisualEffectsDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_visual_effects_disable",
		NameKey:        "booster.system.visual_effects_disable.name",
		DescriptionKey: "booster.system.visual_effects_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "visuals", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.visual_effects_disable.name":        "Отключить визуальные эффекты",
			"booster.system.visual_effects_disable.description": "Сосредоточиться на оптимизации общей среды ОС для повышения эффективности и скорости отклика (делает ваш Windows визуально менее привлекательным).",
		},
		i18n.Spanish: {
			"booster.system.visual_effects_disable.name":        "Desactivar Efectos Visuales",
			"booster.system.visual_effects_disable.description": "Se centra en optimizar el entorno general del sistema operativo para mejorar la eficiencia y la velocidad de respuesta (hace que su Windows sea visualmente más feo).",
		},
		i18n.Portuguese: {
			"booster.system.visual_effects_disable.name":        "Desativar Efeitos Visuais",
			"booster.system.visual_effects_disable.description": "Foca em otimizar o ambiente geral do sistema operativo para melhorar eficiência e velocidade de resposta (Desativa seu Windows visualmente mais feio).",
		},
		i18n.PortugueseBrazil: {
			"booster.system.visual_effects_disable.name":        "Desativar Efeitos Visuais",
			"booster.system.visual_effects_disable.description": "Foca em otimizar o ambiente geral do sistema operacional para melhorar eficiência e velocidade de resposta (Desativa seu Windows visualmente mais feio).",
		},
		i18n.English: {
			"booster.system.visual_effects_disable.name":        "Disable Visual Effects",
			"booster.system.visual_effects_disable.description": "Focuses on optimizing the overall OS environment to improve efficiency and response speed (Makes your Windows visually uglier).",
		},
	}

	executor := NewVisualEffectsDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
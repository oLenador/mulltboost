package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewGameDVRDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_game_dvr_disable",
		NameKey:        "booster.system.game_dvr_disable.name",
		DescriptionKey: "booster.system.game_dvr_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "gaming", "resources"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.game_dvr_disable.name":        "Отключить Game DVR",
			"booster.system.game_dvr_disable.description": "Отключает встроенную запись видео, освобождая ресурсы GPU и ЦП.",
		},
		i18n.Spanish: {
			"booster.system.game_dvr_disable.name":        "Desactivar Game DVR",
			"booster.system.game_dvr_disable.description": "Desactiva la captura de vídeo nativa, ahorrando recursos de GPU y CPU.",
		},
		i18n.Portuguese: {
			"booster.system.game_dvr_disable.name":        "Desativar Game DVR",
			"booster.system.game_dvr_disable.description": "Desliga captura de vídeo nativa, poupando recursos de GPU e CPU.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.game_dvr_disable.name":        "Desativar Game DVR",
			"booster.system.game_dvr_disable.description": "Desliga captura de vídeo nativa, poupando recursos de GPU e CPU.",
		},
		i18n.English: {
			"booster.system.game_dvr_disable.name":        "Disable Game DVR",
			"booster.system.game_dvr_disable.description": "Turns off native video capture, saving GPU and CPU resources.",
		},
	}

	executor := NewGameDVRDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
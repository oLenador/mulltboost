package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewJumboFramesBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_jumbo_frames",
		NameKey:        "booster.connection.jumbo_frames.name",
		DescriptionKey: "booster.connection.jumbo_frames.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "mtu", "jumbo"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.jumbo_frames.name":        "Настроить Jumbo Frames",
			"booster.connection.jumbo_frames.description": "Включает поддержку Jumbo Frames (например, 9000 байт) для уменьшения накладных расходов в гигабитных сетях.",
		},
		i18n.Spanish: {
			"booster.connection.jumbo_frames.name":        "Configurar Jumbo Frames",
			"booster.connection.jumbo_frames.description": "Activa el soporte a Jumbo Frames, permitiendo paquetes más grandes (ej.: 9000 bytes) para reducir el overhead en redes Gigabit.",
		},
		i18n.Portuguese: {
			"booster.connection.jumbo_frames.name":        "Configurar Jumbo Frames",
			"booster.connection.jumbo_frames.description": "Ativa o suporte a Jumbo Frames, permitindo pacotes maiores (ex.: 9000 bytes) para reduzir overhead em redes Gigabit.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.jumbo_frames.name":        "Configurar Jumbo Frames",
			"booster.connection.jumbo_frames.description": "Ativa o suporte a Jumbo Frames, permitindo pacotes maiores (ex.: 9000 bytes) para reduzir overhead em redes Gigabit.",
		},
		i18n.English: {
			"booster.connection.jumbo_frames.name":        "Configure Jumbo Frames",
			"booster.connection.jumbo_frames.description": "Enables jumbo frames support, allowing larger packets (e.g., 9000 bytes) to reduce overhead on Gigabit networks.",
		},
	}

	executor := NewJumboFramesExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
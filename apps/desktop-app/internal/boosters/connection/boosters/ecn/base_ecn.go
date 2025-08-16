package connection

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewECNBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_ecn",
		NameKey:        "booster.connection.ecn.name",
		DescriptionKey: "booster.connection.ecn.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "latency", "ecn"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.ecn.name":        "Включить ECN",
			"booster.connection.ecn.description": "Включает ECN, чтобы сеть могла сигнализировать о перегрузке до отбрасывания пакетов, уменьшая задержку и потери.",
		},
		i18n.Spanish: {
			"booster.connection.ecn.name":        "Habilitar ECN (Explicit Congestion Notification)",
			"booster.connection.ecn.description": "Activa ECN para permitir que la red señale la congestión antes de descartar paquetes, reduciendo la latencia y las pérdidas en conexiones congestionadas.",
		},
		i18n.Portuguese: {
			"booster.connection.ecn.name":        "Habilitar ECN (Explicit Congestion Notification)",
			"booster.connection.ecn.description": "Ativa o ECN para permitir que a rede sinalize congestionamento antes de descartar pacotes, reduzindo latência e perdas em conexões congestionadas.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.ecn.name":        "Habilitar ECN (Explicit Congestion Notification)",
			"booster.connection.ecn.description": "Ativa o ECN para permitir que a rede sinalize congestionamento antes de descartar pacotes, reduzindo latência e perdas em conexões congestionadas.",
		},
		i18n.English: {
			"booster.connection.ecn.name":        "Enable ECN (Explicit Congestion Notification)",
			"booster.connection.ecn.description": "Enables ECN to allow the network to signal congestion before dropping packets, reducing latency and loss on congested connections.",
		},
	}

	executor := NewECNExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
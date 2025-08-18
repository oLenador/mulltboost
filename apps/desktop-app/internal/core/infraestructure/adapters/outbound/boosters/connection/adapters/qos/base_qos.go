package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewQoSBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_qos",
		NameKey:        "booster.connection.qos.name",
		DescriptionKey: "booster.connection.qos.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "qos", "priority"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.qos.name":        "Приоритет пакетов QoS",
			"booster.connection.qos.description": "Приоритезирует важные действия, такие как игры и потоковое видео, уменьшая задержки и улучшая качество соединения.",
		},
		i18n.Spanish: {
			"booster.connection.qos.name":        "Prioridad de Paquetes QoS",
			"booster.connection.qos.description": "Da prioridad a las actividades más importantes, como juegos, streaming y descargas, reduciendo retrasos y mejorando la calidad de la conexión.",
		},
		i18n.Portuguese: {
			"booster.connection.qos.name":        "Prioridade de Pacotes QoS",
			"booster.connection.qos.description": "Dá prioridade às atividades mais importantes, como jogos, streaming e downloads, reduzindo atrasos e melhorando a qualidade da conexão.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.qos.name":        "Prioridade de Pacotes QoS",
			"booster.connection.qos.description": "Dá prioridade às atividades mais importantes, como jogos, streaming e downloads, reduzindo atrasos e melhorando a qualidade da conexão.",
		},
		i18n.English: {
			"booster.connection.qos.name":        "Packet Priority QoS",
			"booster.connection.qos.description": "Prioritizes important activities like gaming, streaming, and downloads, reducing delays and improving connection quality.",
		},
	}

	executor := NewQoSExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
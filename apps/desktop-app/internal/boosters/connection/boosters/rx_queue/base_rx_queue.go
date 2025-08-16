package connection

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewRxQueueBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_rx_queue",
		NameKey:        "booster.connection.rx_queue.name",
		DescriptionKey: "booster.connection.rx_queue.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "queue", "rx"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.rx_queue.name":        "Увеличить размер очереди приема",
			"booster.connection.rx_queue.description": "Увеличивает размер буфера приема TCP для повышения производительности в сетях с высокой задержкой или переменной пропускной способностью.",
		},
		i18n.Spanish: {
			"booster.connection.rx_queue.name":        "Aumentar el Tamaño de la Fila de Recepción",
			"booster.connection.rx_queue.description": "Aumenta el tamaño del buffer de recepción TCP para un mejor rendimiento en redes con alta latencia o ancho de banda variable.",
		},
		i18n.Portuguese: {
			"booster.connection.rx_queue.name":        "Aumentar o Tamanho da Fila de Recebimento",
			"booster.connection.rx_queue.description": "Aumenta o tamanho do buffer de recebimento TCP para melhor desempenho em redes com alta latência ou largura de banda variável.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.rx_queue.name":        "Aumenta o Tamanho da Fila de Recebimento",
			"booster.connection.rx_queue.description": "Aumenta o tamanho do buffer de recebimento TCP para melhor desempenho em redes com alta latência ou largura de banda variável.",
		},
		i18n.English: {
			"booster.connection.rx_queue.name":        "Increase Receive Queue Size",
			"booster.connection.rx_queue.description": "Increases the TCP receive buffer size for better performance on networks with high latency or variable bandwidth.",
		},
	}

	executor := NewRxQueueExecutor() 
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
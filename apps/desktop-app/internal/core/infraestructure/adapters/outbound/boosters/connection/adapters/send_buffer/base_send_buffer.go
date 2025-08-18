package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewSendBufferBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_send_buffer",
		NameKey:        "booster.connection.send_buffer.name",
		DescriptionKey: "booster.connection.send_buffer.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "buffer", "send"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.send_buffer.name":        "Настроить размер буфера отправки",
			"booster.connection.send_buffer.description": "Увеличивает буфер отправки TCP для повышения производительности в сетях с высокой пропускной способностью.",
		},
		i18n.Spanish: {
			"booster.connection.send_buffer.name":        "Ajustar el Tamaño del Send Buffer",
			"booster.connection.send_buffer.description": "Aumenta el buffer de envío TCP para mejorar el rendimiento en redes con un alto ancho de banda.",
		},
		i18n.Portuguese: {
			"booster.connection.send_buffer.name":        "Ajustar o Tamanho do Send Buffer",
			"booster.connection.send_buffer.description": "Aumenta o buffer de envio TCP para melhorar desempenho em redes com alta largura de banda.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.send_buffer.name":        "Ajustar o Tamanho do Send Buffer",
			"booster.connection.send_buffer.description": "Aumenta o buffer de envio TCP para melhorar desempenho em redes com alta largura de banda.",
		},
		i18n.English: {
			"booster.connection.send_buffer.name":        "Adjust Send Buffer Size",
			"booster.connection.send_buffer.description": "Increases the TCP send buffer to improve performance on high-bandwidth networks.",
		},
	}

	executor := NewSendBufferExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
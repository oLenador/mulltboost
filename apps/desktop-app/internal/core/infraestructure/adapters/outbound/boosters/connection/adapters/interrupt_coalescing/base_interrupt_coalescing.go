package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewInterruptCoalescingBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_interrupt_coalescing",
		NameKey:        "booster.connection.interrupt_coalescing.name",
		DescriptionKey: "booster.connection.interrupt_coalescing.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "latency", "coalescing"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.interrupt_coalescing.name":        "Отключить объединение прерываний",
			"booster.connection.interrupt_coalescing.description": "Отключает объединение прерываний на сетевой карте, уменьшая задержку за счет более высокого использования ЦП.",
		},
		i18n.Spanish: {
			"booster.connection.interrupt_coalescing.name":        "Desactivar Interrupt Coalescing",
			"booster.connection.interrupt_coalescing.description": "Desactiva la coalescencia de interrupciones en la tarjeta de red, reduciendo el retraso en la entrega de paquetes al costo de un mayor uso de la CPU.",
		},
		i18n.Portuguese: {
			"booster.connection.interrupt_coalescing.name":        "Desativar Interrupt Coalescing",
			"booster.connection.interrupt_coalescing.description": "Desativa a coalescência de interrupções na placa de rede, reduzindo o atraso na entrega de pacotes ao custo de maior uso da CPU.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.interrupt_coalescing.name":        "Desativar Interrupt Coalescing",
			"booster.connection.interrupt_coalescing.description": "Desativa a coalescência de interrupções na placa de rede, reduzindo o atraso na entrega de pacotes ao custo de maior uso da CPU.",
		},
		i18n.English: {
			"booster.connection.interrupt_coalescing.name":        "Disable Interrupt Coalescing",
			"booster.connection.interrupt_coalescing.description": "Disables interrupt coalescing on the network card, reducing packet delivery latency at the cost of higher CPU usage.",
		},
	}

	executor := NewInterruptCoalescingExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
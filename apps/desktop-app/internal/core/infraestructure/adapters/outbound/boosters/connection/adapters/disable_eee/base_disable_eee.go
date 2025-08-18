package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewEEEBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_eee",
		NameKey:        "booster.connection.eee.name",
		DescriptionKey: "booster.connection.eee.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "energy", "latency"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.eee.name":        "Отключить Energy-Efficient Ethernet",
			"booster.connection.eee.description": "Отключает функцию энергосбережения в сетевой карте, которая может вносить задержку в обмен на энергоэффективность.",
		},
		i18n.Spanish: {
			"booster.connection.eee.name":        "Desactivar Energy-Efficient Ethernet",
			"booster.connection.eee.description": "Desactiva el recurso de ahorro de energía en la tarjeta de red (EEE), que puede introducir latencia a cambio de eficiencia energética.",
		},
		i18n.Portuguese: {
			"booster.connection.eee.name":        "Desativar Energy-Efficient Ethernet",
			"booster.connection.eee.description": "Desativa o recurso de poupança de energia na placa de rede (EEE), que pode introduzir latência em troca de eficiência energética.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.eee.name":        "Desativar Energy-Efficient Ethernet",
			"booster.connection.eee.description": "Desativa o recurso de economia de energia na placa de rede (EEE), que pode introduzir latência em troca de eficiência energética.",
		},
		i18n.English: {
			"booster.connection.eee.name":        "Disable Energy-Efficient Ethernet",
			"booster.connection.eee.description": "Disables the energy-saving feature on the network card (EEE), which may introduce latency in exchange for energy efficiency.",
		},
	}

	executor := NewEEEExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
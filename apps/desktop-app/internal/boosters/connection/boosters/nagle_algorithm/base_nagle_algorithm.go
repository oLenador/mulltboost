package connection

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewNagleAlgorithmBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_nagle_algorithm",
		NameKey:        "booster.connection.nagle_algorithm.name",
		DescriptionKey: "booster.connection.nagle_algorithm.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "latency", "nagle"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.nagle_algorithm.name":        "Отключить алгоритм Нейгла",
			"booster.connection.nagle_algorithm.description": "Отключает алгоритм Нейгла, который объединяет небольшие пакеты, уменьшая задержку в реальном времени.",
		},
		i18n.Spanish: {
			"booster.connection.nagle_algorithm.name":        "Desactivar el Algoritmo de Nagle",
			"booster.connection.nagle_algorithm.description": "Desactiva el algoritmo de Nagle, que agrupa pequeños paquetes antes de enviarlos, reduciendo la latencia en aplicaciones en tiempo real.",
		},
		i18n.Portuguese: {
			"booster.connection.nagle_algorithm.name":        "Desativar o Algoritmo de Nagle",
			"booster.connection.nagle_algorithm.description": "Desativa o algoritmo de Nagle, que agrupa pequenos pacotes antes de os enviar, reduzindo a latência em aplicações em tempo real.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.nagle_algorithm.name":        "Desativar o Algoritmo de Nagle",
			"booster.connection.nagle_algorithm.description": "Desativa o algoritmo de Nagle, que agrupa pequenos pacotes antes de enviá-los, reduzindo a latência em aplicações em tempo real.",
		},
		i18n.English: {
			"booster.connection.nagle_algorithm.name":        "Disable Nagle's Algorithm",
			"booster.connection.nagle_algorithm.description": "Disables Nagle's algorithm, which groups small packets before sending them, reducing latency in real-time applications.",
		},
	}

	executor := NewNagleAlgorithmExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
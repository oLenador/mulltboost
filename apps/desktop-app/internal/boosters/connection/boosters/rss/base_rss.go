package connection

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewRSSBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_rss",
		NameKey:        "booster.connection.rss.name",
		DescriptionKey: "booster.connection.rss.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "rss", "queues"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.rss.name":        "Включить RSS и многопоточные очереди",
			"booster.connection.rss.description": "Активирует RSS для распределения обработки пакетов между несколькими ядрами ЦП, увеличивая количество очередей для улучшения производительности.",
		},
		i18n.Spanish: {
			"booster.connection.rss.name":        "Habilitar RSS (Receive Side Scaling) con Filas Múltiples",
			"booster.connection.rss.description": "Activa las configuraciones RSS para distribuir el procesamiento de paquetes recibidos entre múltiples núcleos de la CPU, aumentando el número de filas (queues) para mejorar el rendimiento.",
		},
		i18n.Portuguese: {
			"booster.connection.rss.name":        "Habilitar RSS (Receive Side Scaling) com Filas Múltiplas",
			"booster.connection.rss.description": "Ativa configurações RSS para distribuir o processamento de pacotes recebidos entre múltiplos núcleos da CPU, aumentando o número de filas (queues) para melhorar o desempenho.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.rss.name":        "Habilitar RSS (Receive Side Scaling) com Filas Múltiplas",
			"booster.connection.rss.description": "Ativa configurações RSS para distribuir o processamento de pacotes recebidos entre múltiplos núcleos da CPU, aumentando o número de filas (queues) para melhorar o desempenho.",
		},
		i18n.English: {
			"booster.connection.rss.name":        "Enable RSS (Receive Side Scaling) with Multiple Queues",
			"booster.connection.rss.description": "Enables RSS configurations to distribute the processing of received packets among multiple CPU cores, increasing the number of queues for better performance.",
		},
	}

	executor := NewRSSBooster() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
package connection

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewARPCacheBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_arp_cache",
		NameKey:        "booster.connection.arp_cache.name",
		DescriptionKey: "booster.connection.arp_cache.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "arp", "cache"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.arp_cache.name":        "Уменьшить время кэша ARP",
			"booster.connection.arp_cache.description": "Сокращает время жизни кэша ARP для более частых обновлений в динамических сетях.",
		},
		i18n.Spanish: {
			"booster.connection.arp_cache.name":        "Reducir Tiempo de Caché ARP",
			"booster.connection.arp_cache.description": "Disminuye el tiempo de vida de la caché ARP para forzar actualizaciones más frecuentes, útil en redes dinámicas.",
		},
		i18n.Portuguese: {
			"booster.connection.arp_cache.name":        "Reduzir Tempo de Cache ARP",
			"booster.connection.arp_cache.description": "Diminui o tempo de vida do cache ARP para forçar atualizações mais frequentes, útil em redes dinâmicas.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.arp_cache.name":        "Reduzir Tempo de Cache ARP",
			"booster.connection.arp_cache.description": "Diminui o tempo de vida do cache ARP para forçar atualizações mais frequentes, útil em redes dinâmicas.",
		},
		i18n.English: {
			"booster.connection.arp_cache.name":        "Reduce ARP Cache Timeout",
			"booster.connection.arp_cache.description": "Decreases the ARP cache lifetime to force more frequent updates, useful in dynamic networks.",
		},
	}

	executor := NewARPCacheExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
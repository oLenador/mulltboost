package connection

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewDNSCacheBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_dns_cache",
		NameKey:        "booster.connection.dns_cache.name",
		DescriptionKey: "booster.connection.dns_cache.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "dns", "cache"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.dns_cache.name":        "Очистка кэша DNS",
			"booster.connection.dns_cache.description": "Очищает кэш DNS для устранения проблем с подключением и обеспечения использования последних IP-адресов.",
		},
		i18n.Spanish: {
			"booster.connection.dns_cache.name":        "Limpieza de Caché DNS",
			"booster.connection.dns_cache.description": "Limpia la caché DNS para resolver problemas de conectividad y asegurar que el sistema use la información de IP más reciente.",
		},
		i18n.Portuguese: {
			"booster.connection.dns_cache.name":        "Limpeza de Cache DNS",
			"booster.connection.dns_cache.description": "Limpa o cache DNS para resolver problemas de conectividade e garantir que o sistema está a usar as informações de IP mais recentes.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.dns_cache.name":        "Limpeza de Cache DNS",
			"booster.connection.dns_cache.description": "Limpa o cache DNS para resolver problemas de conectividade e garantir que o sistema está usando as informações de IP mais recentes.",
		},
		i18n.English: {
			"booster.connection.dns_cache.name":        "Clear DNS Cache",
			"booster.connection.dns_cache.description": "Clears the DNS cache to resolve connectivity issues and ensure the system is using the most recent IP information.",
		},
	}

	executor := NewDNSCacheExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
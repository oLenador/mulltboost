package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

type DNSBooster struct {
	*booster.BaseBooster
}

func NewDNSBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_dns_booster",
		NameKey:        "booster.connection.dns.name",
		DescriptionKey: "booster.connection.dns.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "dns", "speed"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.English: {
			"booster.connection.dns.name":        "DNS Optimizer",
			"booster.connection.dns.description": "Configures faster DNS servers and clears DNS cache",
		},
		i18n.Portuguese: {
			"booster.connection.dns.name":        "Otimizador DNS",
			"booster.connection.dns.description": "Configura servidores DNS mais rápidos e limpa o cache DNS",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.dns.name":        "Otimizador de DNS",
			"booster.connection.dns.description": "Configura servidores de DNS mais rápidos e limpa o cache do DNS",
		},
		i18n.Spanish: {
			"booster.connection.dns.name":        "Optimizador DNS",
			"booster.connection.dns.description": "Configura servidores DNS más rápidos y limpia la caché DNS",
		},
		i18n.Russian: {
			"booster.connection.dns.name":        "Оптимизатор DNS",
			"booster.connection.dns.description": "Настраивает более быстрые DNS-серверы и очищает кеш DNS",
		},
	}	

	i18nSvc := i18n.NewService()
	i18nSvc.SetTranslations(translations)


	
	executor := NewDNSCacheExecutor(services)
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	

	return baseBooster
}
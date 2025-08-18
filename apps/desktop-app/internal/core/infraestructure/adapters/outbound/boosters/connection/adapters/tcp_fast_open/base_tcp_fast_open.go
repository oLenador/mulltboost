//go:build windows
package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewTCPFastOpenBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_tcp_fast_open",
		NameKey:        "booster.connection.tcp_fast_open.name",
		DescriptionKey: "booster.connection.tcp_fast_open.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "tcp", "fast open"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.tcp_fast_open.name":        "Включить TCP Fast Open",
			"booster.connection.tcp_fast_open.description": "Включает TCP Fast Open, что сокращает время начального рукопожатия TCP-соединений, улучшая начальную задержку.",
		},
		i18n.Spanish: {
			"booster.connection.tcp_fast_open.name":        "Habilitar TCP Fast Open",
			"booster.connection.tcp_fast_open.description": "Activa el TCP Fast Open, que reduce el tiempo de handshake en conexiones TCP, mejorando la latencia inicial.",
		},
		i18n.Portuguese: {
			"booster.connection.tcp_fast_open.name":        "Habilitar TCP Fast Open",
			"booster.connection.tcp_fast_open.description": "Ativa o TCP Fast Open, que reduz o tempo de handshake em conexões TCP, melhorando a latência inicial.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.tcp_fast_open.name":        "Habilitar TCP Fast Open",
			"booster.connection.tcp_fast_open.description": "Ativa o TCP Fast Open, que reduz o tempo de handshake em conexões TCP, melhorando a latência inicial.",
		},
		i18n.English: {
			"booster.connection.tcp_fast_open.name":        "Enable TCP Fast Open",
			"booster.connection.tcp_fast_open.description": "Enables TCP Fast Open, which reduces the handshake time on TCP connections, improving initial latency.",
		},
	}

	executor := NewTCPFastOpenExecutor(
		services.RegistryService,
		services.SystemService,
		services.ElevationService,
	) 
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
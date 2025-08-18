//go:build windows
package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewTCPRTOBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_tcp_rto",
		NameKey:        "booster.connection.tcp_rto.name",
		DescriptionKey: "booster.connection.tcp_rto.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "tcp", "rto"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.tcp_rto.name":        "Настроить интервал TCP-ретрансляции",
			"booster.connection.tcp_rto.description": "Уменьшает начальный интервал ретрансляции TCP (RTO) для ускорения восстановления потерянных пакетов.",
		},
		i18n.Spanish: {
			"booster.connection.tcp_rto.name":        "Ajustar el Intervalo de Retransmisión TCP",
			"booster.connection.tcp_rto.description": "Reduce el intervalo inicial de retransmisión TCP (RTO) para acelerar la recuperación de paquetes perdidos.",
		},
		i18n.Portuguese: {
			"booster.connection.tcp_rto.name":        "Ajustar o Intervalo de Retransmissão TCP",
			"booster.connection.tcp_rto.description": "Reduz o intervalo inicial de retransmissão TCP (RTO) para acelerar a recuperação de pacotes perdidos.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.tcp_rto.name":        "Ajustar o Intervalo de Retransmissão TCP",
			"booster.connection.tcp_rto.description": "Reduz o intervalo inicial de retransmissão TCP (RTO) para acelerar a recuperação de pacotes perdidos.",
		},
		i18n.English: {
			"booster.connection.tcp_rto.name":        "Adjust TCP Retransmission Interval",
			"booster.connection.tcp_rto.description": "Reduces the initial TCP retransmission interval (RTO) to speed up recovery from lost packets.",
		},
	}

	executor := NewTCPRTOExecutor(
		services.RegistryService,
		services.SystemService,
		services.ElevationService,
	)
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
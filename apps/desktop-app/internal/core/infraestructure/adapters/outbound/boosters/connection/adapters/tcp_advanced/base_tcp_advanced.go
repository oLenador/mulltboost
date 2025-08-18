//go:build windows
package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewAdvancedTCPBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_tcp_advanced",
		NameKey:        "booster.connection.tcp_advanced.name",
		DescriptionKey: "booster.connection.tcp_advanced.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "tcp", "advanced"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.tcp_advanced.name":        "Продвинутая оптимизация TCP/IP",
			"booster.connection.tcp_advanced.description": "Настраивает продвинутые параметры TCP/IP для ускорения одновременных задач, идеально для быстрой загрузки и плавной навигации.",
		},
		i18n.Spanish: {
			"booster.connection.tcp_advanced.name":        "Optimización Avanzada de Configuraciones TCP/IP",
			"booster.connection.tcp_advanced.description": "Ajusta parámetros para acelerar más tareas al mismo tiempo, ideal para descargas rápidas y navegación sin interrupciones.",
		},
		i18n.Portuguese: {
			"booster.connection.tcp_advanced.name":        "Otimização Avançada de Configurações TCP/IP",
			"booster.connection.tcp_advanced.description": "Ajusta parâmetros para acelerar mais tarefas ao mesmo tempo, ideal para downloads rápidos e navegação sem interrupções.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.tcp_advanced.name":        "Otimização Avançada de Configurações TCP/IP",
			"booster.connection.tcp_advanced.description": "Ajusta parâmetros para acelerar mais tarefas ao mesmo tempo, ideal para downloads rápidos e navegação sem interrupções.",
		},
		i18n.English: {
			"booster.connection.tcp_advanced.name":        "Advanced TCP/IP Configuration Optimization",
			"booster.connection.tcp_advanced.description": "Adjusts parameters to accelerate multiple tasks simultaneously, ideal for fast downloads and uninterrupted browsing.",
		},
	}

	executor := NewTCPAdvancedExecutor(
		services.TcpService,
		services.RegistryService,
		services.SystemService,
		services.ElevationService,
	)
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
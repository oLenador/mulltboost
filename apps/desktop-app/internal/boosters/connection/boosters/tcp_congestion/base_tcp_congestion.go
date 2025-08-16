package connection

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewTCPCongestionBooster() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_tcp_congestion",
		NameKey:        "booster.connection.tcp_congestion.name",
		DescriptionKey: "booster.connection.tcp_congestion.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelPremium,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "tcp", "congestion"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.tcp_congestion.name":        "Настроить TCP-перегрузку (CTCP или Cubic)",
			"booster.connection.tcp_congestion.description": "Настраивает алгоритм управления TCP-перегрузкой для использования более агрессивных методов, таких как CTCP или Cubic, оптимизируя пропускную способность.",
		},
		i18n.Spanish: {
			"booster.connection.tcp_congestion.name":        "Ajustar el Congestionamiento TCP (CTCP o Cubic)",
			"booster.connection.tcp_congestion.description": "Configura el algoritmo de control de congestionamiento TCP para usar más agresivo como CTCP (Compound TCP) o Cubic, optimizando el rendimiento en conexiones de alto ancho de banda.",
		},
		i18n.Portuguese: {
			"booster.connection.tcp_congestion.name":        "Ajustar o Congestionamento TCP (CTCP ou Cubic)",
			"booster.connection.tcp_congestion.description": "Configura o algoritmo de controlo de congestionamento TCP para usar mais agressivo como CTCP (Compound TCP) ou Cubic, otimizando o throughput em conexões de alta largura de banda.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.tcp_congestion.name":        "Ajustar o Congestionamento TCP (CTCP ou Cubic)",
			"booster.connection.tcp_congestion.description": "Configura o algoritmo de controle de congestionamento TCP para usar mais agressivo como CTCP (Compound TCP) ou Cubic, otimizando o throughput em conexões de alta largura de banda.",
		},
		i18n.English: {
			"booster.connection.tcp_congestion.name":        "Adjust TCP Congestion (CTCP or Cubic)",
			"booster.connection.tcp_congestion.description": "Configures the TCP congestion control algorithm to use more aggressive methods like CTCP (Compound TCP) or Cubic, optimizing throughput on high-bandwidth connections.",
		},
	}

	executor := NewTCPCongestionExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
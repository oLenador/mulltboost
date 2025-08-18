package connection

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"

	"github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewMTUAdjustmentBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "connection_mtu_adjustment",
		NameKey:        "booster.connection.mtu_adjustment.name",
		DescriptionKey: "booster.connection.mtu_adjustment.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "mtu", "optimization"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.connection.mtu_adjustment.name":        "Автоматическая настройка MTU",
			"booster.connection.mtu_adjustment.description": "Настраивает оптимальный размер пакетов данных, чтобы избежать фрагментации и улучшить стабильность соединения.",
		},
		i18n.Spanish: {
			"booster.connection.mtu_adjustment.name":        "Ajuste automático de MTU",
			"booster.connection.mtu_adjustment.description": "Ajusta el tamaño ideal de los paquetes de datos para evitar la fragmentación y mejorar la estabilidad y velocidad de la conexión.",
		},
		i18n.Portuguese: {
			"booster.connection.mtu_adjustment.name":        "Ajuste automático de MTU",
			"booster.connection.mtu_adjustment.description": "Ajusta o tamanho ideal dos pacotes de dados para evitar fragmentação e melhorar a estabilidade e velocidade da ligação.",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.mtu_adjustment.name":        "Ajuste automático de MTU",
			"booster.connection.mtu_adjustment.description": "Ajusta o tamanho ideal dos pacotes de dados para evitar fragmentação e melhorar a estabilidade e velocidade da conexão.",
		},
		i18n.English: {
			"booster.connection.mtu_adjustment.name":        "Automatic MTU Adjustment",
			"booster.connection.mtu_adjustment.description": "Adjusts the ideal data packet size to avoid fragmentation and improve connection stability and speed.",
		},
	}

	executor := NewMTUAdjustmentExecutor() // Adicione a implementação real do executor
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
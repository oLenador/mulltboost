package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewTRIMSSD() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_trim_ssd",
		NameKey:        "booster.system.trim_ssd.name",
		DescriptionKey: "booster.system.trim_ssd.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "ssd", "trim"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.trim_ssd.name":        "Включить TRIM на SSD",
			"booster.system.trim_ssd.description": "Активирует автоматическую очистку на SSD, поддерживая скорость чтения и записи.",
		},
		i18n.Spanish: {
			"booster.system.trim_ssd.name":        "Habilitar TRIM en SSD",
			"booster.system.trim_ssd.description": "Activa la limpieza automática en el SSD, manteniendo la velocidad en lecturas y escrituras.",
		},
		i18n.Portuguese: {
			"booster.system.trim_ssd.name":        "Habilitar TRIM em SSD",
			"booster.system.trim_ssd.description": "Ativa limpeza automática no SSD, mantendo velocidade em leituras e gravações.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.trim_ssd.name":        "Habilitar TRIM em SSD",
			"booster.system.trim_ssd.description": "Ativa limpeza automática no SSD, mantendo velocidade em leituras e gravações.",
		},
		i18n.English: {
			"booster.system.trim_ssd.name":        "Enable TRIM on SSD",
			"booster.system.trim_ssd.description": "Enables automatic cleanup on the SSD, maintaining read and write speeds.",
		},
	}

	executor := NewTRIMSSDExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}

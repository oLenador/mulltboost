package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewTelemetryDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_telemetry_disable",
		NameKey:        "booster.system.telemetry_disable.name",
		DescriptionKey: "booster.system.telemetry_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "telemetry", "privacy"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.telemetry_disable.name":        "Отключить телеметрию",
			"booster.system.telemetry_disable.description": "Прекращает отправку диагностических данных, экономя ресурсы ЦП и пропускную способность.",
		},
		i18n.Spanish: {
			"booster.system.telemetry_disable.name":        "Desactivar Telemetría",
			"booster.system.telemetry_disable.description": "Interrumpe el envío de diagnósticos, ahorrando CPU y ancho de banda del sistema.",
		},
		i18n.Portuguese: {
			"booster.system.telemetry_disable.name":        "Desativar Telemetria",
			"booster.system.telemetry_disable.description": "Interrompe envio de diagnósticos, economizando CPU e largura de banda do sistema.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.telemetry_disable.name":        "Desativar Telemetria",
			"booster.system.telemetry_disable.description": "Interrompe envio de diagnósticos, economizando CPU e largura de banda do sistema.",
		},
		i18n.English: {
			"booster.system.telemetry_disable.name":        "Disable Telemetry",
			"booster.system.telemetry_disable.description": "Stops sending diagnostic data, saving CPU and system bandwidth.",
		},
	}

	executor := NewTelemetryDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
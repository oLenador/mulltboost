package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewErrorReportingDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_error_reporting_disable",
		NameKey:        "booster.system.error_reporting_disable.name",
		DescriptionKey: "booster.system.error_reporting_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "errors", "resources"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.error_reporting_disable.name":        "Отключить отчеты об ошибках",
			"booster.system.error_reporting_disable.description": "Останавливает отправку отчетов об ошибках, экономя ресурсы ЦП и диска.",
		},
		i18n.Spanish: {
			"booster.system.error_reporting_disable.name":        "Desactivar el Reporte de Errores",
			"booster.system.error_reporting_disable.description": "Detiene el envío de informes de fallas, ahorrando CPU y disco en errores.",
		},
		i18n.Portuguese: {
			"booster.system.error_reporting_disable.name":        "Desligar Relatório de Erros",
			"booster.system.error_reporting_disable.description": "Para o envio de relatórios de falhas, economizando CPU e disco em erros.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.error_reporting_disable.name":        "Desligar Relatório de Erros",
			"booster.system.error_reporting_disable.description": "Para o envio de relatórios de falhas, economizando CPU e disco em erros.",
		},
		i18n.English: {
			"booster.system.error_reporting_disable.name":        "Turn Off Error Reporting",
			"booster.system.error_reporting_disable.description": "Stops sending crash reports, saving CPU and disk in case of errors.",
		},
	}

	executor := NewErrorReportingDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
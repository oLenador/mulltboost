package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewAppBackgroundDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_app_background_disable",
		NameKey:        "booster.system.app_background_disable.name",
		DescriptionKey: "booster.system.app_background_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "performance", "ram"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.app_background_disable.name":        "Отключить фоновые приложения",
			"booster.system.app_background_disable.description": "Предотвращает работу приложений в фоновом режиме, освобождая ЦП и ОЗУ для приоритетных задач.",
		},
		i18n.Spanish: {
			"booster.system.app_background_disable.name":        "Desactivar Aplicaciones en Segundo Plano",
			"booster.system.app_background_disable.description": "Impide que las aplicaciones se ejecuten en segundo plano, liberando CPU y RAM para tareas prioritarias.",
		},
		i18n.Portuguese: {
			"booster.system.app_background_disable.name":        "Desativar Aplicações em Segundo Plano",
			"booster.system.app_background_disable.description": "Impede aplicações em segundo plano, libertando CPU e RAM para tarefas prioritárias.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.app_background_disable.name":        "Desativar Background Apps",
			"booster.system.app_background_disable.description": "Impede aplicativos em segundo plano, liberando CPU e RAM para tarefas prioritárias.",
		},
		i18n.English: {
			"booster.system.app_background_disable.name":        "Disable Background Apps",
			"booster.system.app_background_disable.description": "Prevents applications from running in the background, freeing up CPU and RAM for priority tasks.",
		},
	}

	executor := NewAppBackgroundDisableExecutor() // Add the real executor implementation
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
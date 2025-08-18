package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewSysMainOptimize() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_sysmain_optimize",
		NameKey:        "booster.system.sysmain_optimize.name",
		DescriptionKey: "booster.system.sysmain_optimize.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "memory", "startup"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.sysmain_optimize.name":        "Оптимизировать SysMain",
			"booster.system.sysmain_optimize.description": "Предварительно загружает критические файлы в ОЗУ, ускоряя запуск и доступ к данным.",
		},
		i18n.Spanish: {
			"booster.system.sysmain_optimize.name":        "Optimizar SysMain",
			"booster.system.sysmain_optimize.description": "Precarga archivos críticos en la RAM, acelerando el inicio y el acceso a los datos.",
		},
		i18n.Portuguese: {
			"booster.system.sysmain_optimize.name":        "Otimizar SysMain",
			"booster.system.sysmain_optimize.description": "Pré-carrega ficheiros críticos em RAM, acelerando inicialização e acesso a dados.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.sysmain_optimize.name":        "Otimizar SysMain",
			"booster.system.sysmain_optimize.description": "Pré-carrega arquivos críticos em RAM, acelerando inicialização e acesso a dados.",
		},
		i18n.English: {
			"booster.system.sysmain_optimize.name":        "Optimize SysMain",
			"booster.system.sysmain_optimize.description": "Pre-loads critical files into RAM, speeding up startup and data access.",
		},
	}

	executor := NewSysMainOptimizeExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
package memory

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewPageFileSize() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "memory_page_file_size",
		NameKey:        "booster.memory.page_file_size.name",
		DescriptionKey: "booster.memory.page_file_size.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"memory", "pagefile", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.memory.page_file_size.name":        "Настроить размер файла подкачки",
			"booster.memory.page_file_size.description": "Устанавливает фиксированный размер файла подкачки, предотвращая динамическое выделение и задержку при заполнении ОЗУ.",
		},
		i18n.Spanish: {
			"booster.memory.page_file_size.name":        "Ajustar el Tamaño del Archivo de Paginación",
			"booster.memory.page_file_size.description": "Define un archivo de paginación fijo, evitando la asignación dinámica y la latencia cuando la RAM está llena.",
		},
		i18n.Portuguese: {
			"booster.memory.page_file_size.name":        "Ajustar Tamanho do Arquivo de Paginação",
			"booster.memory.page_file_size.description": "Define arquivo de paginação fixo, evitando alocação dinâmica e latência em RAM cheia.",
		},
		i18n.PortugueseBrazil: {
			"booster.memory.page_file_size.name":        "Ajustar Tamanho do Arquivo de Paginação",
			"booster.memory.page_file_size.description": "Define arquivo de paginação fixo, evitando alocação dinâmica e latência em RAM cheia.",
		},
		i18n.English: {
			"booster.memory.page_file_size.name":        "Adjust Page File Size",
			"booster.memory.page_file_size.description": "Sets a fixed page file size, avoiding dynamic allocation and latency when RAM is full.",
		},
	}

	executor := NewPageFileSizeExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
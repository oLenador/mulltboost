package memory

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewMemoryFree() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "memory_free",
		NameKey:        "booster.memory.free.name",
		DescriptionKey: "booster.memory.free.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     false,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"memory", "cleanup", "performance"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.memory.free.name":        "Освободить память",
			"booster.memory.free.description": "Очищает ненужные файлы и электронный мусор, что мешает производительности вашего компьютера.",
		},
		i18n.Spanish: {
			"booster.memory.free.name":        "Liberar Memoria",
			"booster.memory.free.description": "Limpia archivos innecesarios y basura electrónica que perjudican el rendimiento de su computadora.",
		},
		i18n.Portuguese: {
			"booster.memory.free.name":        "Liberar Memória",
			"booster.memory.free.description": "Limpa ficheiros desnecessários e lixo eletrónico que prejudicam o desempenho do seu computador.",
		},
		i18n.PortugueseBrazil: {
			"booster.memory.free.name":        "Liberar Memória",
			"booster.memory.free.description": "Limpa arquivos desnecessários e lixo eletrônico que prejudicam o desempenho do seu computador.",
		},
		i18n.English: {
			"booster.memory.free.name":        "Free Up Memory",
			"booster.memory.free.description": "Cleans up unnecessary files and junk that harm your computer's performance.",
		},
	}

	executor := NewMemoryFreeExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
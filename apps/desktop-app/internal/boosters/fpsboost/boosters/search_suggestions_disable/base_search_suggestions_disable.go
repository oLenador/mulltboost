package system

import (
	booster "github.com/oLenador/mulltbost/internal/boosters/base"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func NewSearchSuggestionsDisable() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_search_suggestions_disable",
		NameKey:        "booster.system.search_suggestions_disable.name",
		DescriptionKey: "booster.system.search_suggestions_disable.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "search", "privacy"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.search_suggestions_disable.name":        "Отключить поисковые предложения",
			"booster.system.search_suggestions_disable.description": "Удаляет системные подсказки поиска, сокращая обработку и сетевой трафик.",
		},
		i18n.Spanish: {
			"booster.system.search_suggestions_disable.name":        "Desactivar Sugerencias de Búsqueda",
			"booster.system.search_suggestions_disable.description": "Elimina las sugerencias de búsqueda del sistema, cortando el procesamiento y el tráfico de red.",
		},
		i18n.Portuguese: {
			"booster.system.search_suggestions_disable.name":        "Desativar Sugestões de Pesquisa",
			"booster.system.search_suggestions_disable.description": "Remove dicas de busca do sistema, cortando processamento e tráfego de rede.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.search_suggestions_disable.name":        "Desativar Sugestões de Pesquisa",
			"booster.system.search_suggestions_disable.description": "Remove dicas de busca do sistema, cortando processamento e tráfego de rede.",
		},
		i18n.English: {
			"booster.system.search_suggestions_disable.name":        "Disable Search Suggestions",
			"booster.system.search_suggestions_disable.description": "Removes system search suggestions, cutting down on processing and network traffic.",
		},
	}

	executor := NewSearchSuggestionsDisableExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
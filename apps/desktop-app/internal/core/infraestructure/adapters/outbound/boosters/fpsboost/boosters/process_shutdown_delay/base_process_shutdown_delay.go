package system

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)

func NewProcessShutdownDelay() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "system_process_shutdown_delay",
		NameKey:        "booster.system.process_shutdown_delay.name",
		DescriptionKey: "booster.system.process_shutdown_delay.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"system", "processes", "speed"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.system.process_shutdown_delay.name":        "Уменьшить задержку завершения процессов",
			"booster.system.process_shutdown_delay.description": "Сокращает время ожидания для завершения задач, ускоряя освобождение неактивных ресурсов.",
		},
		i18n.Spanish: {
			"booster.system.process_shutdown_delay.name":        "Reducir el Tiempo de Apagado de Procesos",
			"booster.system.process_shutdown_delay.description": "Disminuye la espera para finalizar tareas, acelerando la liberación de recursos ociosos.",
		},
		i18n.Portuguese: {
			"booster.system.process_shutdown_delay.name":        "Reduzir Tempo de Desligamento de Processos",
			"booster.system.process_shutdown_delay.description": "Diminui espera para encerrar tarefas, acelerando liberação de recursos ociosos.",
		},
		i18n.PortugueseBrazil: {
			"booster.system.process_shutdown_delay.name":        "Reduzir Tempo de Desligamento de Processos",
			"booster.system.process_shutdown_delay.description": "Diminui espera para encerrar tarefas, acelerando liberação de recursos ociosos.",
		},
		i18n.English: {
			"booster.system.process_shutdown_delay.name":        "Reduce Process Shutdown Delay",
			"booster.system.process_shutdown_delay.description": "Decreases the wait time to end tasks, speeding up the release of idle resources.",
		},
	}

	executor := NewProcessShutdownDelayExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
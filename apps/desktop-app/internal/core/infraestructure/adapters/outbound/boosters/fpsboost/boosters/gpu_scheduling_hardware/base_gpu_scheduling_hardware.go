package gpu

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	booster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/base"
)


func NewGPUSchedulingHardware() inbound.BoosterUseCase {
	entity := entities.Booster{
		ID:             "gpu_scheduling_hardware",
		NameKey:        "booster.gpu.scheduling_hardware.name",
		DescriptionKey: "booster.gpu.scheduling_hardware.description",
		Category:       entities.CategorySystem,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"gpu", "scheduling", "hardware"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.Russian: {
			"booster.gpu.scheduling_hardware.name":        "Включить аппаратное планирование GPU",
			"booster.gpu.scheduling_hardware.description": "Передает рендеринг GPU, облегчая нагрузку на ЦП при интенсивных графических задачах.",
		},
		i18n.Spanish: {
			"booster.gpu.scheduling_hardware.name":        "Habilitar Agendamiento de GPU con Hardware",
			"booster.gpu.scheduling_hardware.description": "Delega la renderización a la GPU, aliviando la CPU en tareas gráficas intensas.",
		},
		i18n.Portuguese: {
			"booster.gpu.scheduling_hardware.name":        "Habilitar Agendamento de GPU com Hardware",
			"booster.gpu.scheduling_hardware.description": "Delega renderização à GPU, aliviando a CPU em tarefas gráficas intensas.",
		},
		i18n.PortugueseBrazil: {
			"booster.gpu.scheduling_hardware.name":        "Habilitar GPU Scheduling com Hardware",
			"booster.gpu.scheduling_hardware.description": "Delega renderização à GPU, aliviando a CPU em tarefas gráficas intensas.",
		},
		i18n.English: {
			"booster.gpu.scheduling_hardware.name":        "Enable GPU Scheduling with Hardware",
			"booster.gpu.scheduling_hardware.description": "Delegates rendering to the GPU, relieving the CPU from intense graphical tasks.",
		},
	}

	executor := NewGPUSchedulingHardwareExecutor()
	baseBooster := booster.NewBaseBooster(entity, translations, executor)
	return baseBooster
}
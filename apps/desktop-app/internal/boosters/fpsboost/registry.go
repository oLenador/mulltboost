package performance

import (
	appBackgroundDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/app_background_disable"
	bluetoothDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/bluetooth_disable"
	cpuAffinityBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/cpu_affinity"
	cpuGamePriorityBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/cpu_game_priority"
	cpuMitigationsBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/cpu_mitigations"
	cpuSystemServicesReduceBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/cpu_system_services_reduce"
	explorerOptimizeBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/explorer_optimize"
	errorReportingDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/error_reporting_disable"
	fastStartupDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/fast_startup_disable"
	gameDvrDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/game_dvr_disable"
	gpuLowLatencyModeBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/gpu_low_latency_mode"
	gpuSchedulingBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/gpu_scheduling"
	gpuSchedulingHardwareBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/gpu_scheduling_hardware"
	gpuIrqOptimizeBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/gpu_irq_optimize"
	gpuPowerPolicyBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/gpu_power_policy"
	hibernationDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/hibernation_disable"
	mapsManagerDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/maps_manager_disable"
	memoryAllocationBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/memory_allocation"
	memoryCacheOptimizeBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/memory_cache_optimize"
	memoryFlushBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/memory_flush"
	memoryFreeBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/memory_free"
	multBoostPowerBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/mult_boost_power"
	nonPagedPoolOptimizeBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/non_paged_pool_optimize"
	pageFileSizeBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/page_file_size"
	prefetchSuperfetchDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/prefetch_superfetch_disable"
	processShutdownDelayBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/process_shutdown_delay"
	searchSuggestionsDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/search_suggestions_disable"
	servicesDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/services_disable"
	smartScreenDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/smart_screen_disable"
	standbyListCleanBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/standby_list_clean"
	sysmainOptimizeBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/sysmain_optimize"
	telemetryDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/telemetry_disable"
	threadSchedulingBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/thread_scheduling"
	timerResolutionBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/timer_resolution"
	trimSsdBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/trim_ssd"
	visualEffectsDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/visual_effects_disable"
	windowsDefenderDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/windows_defender_disable"
	windowsSearchDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/windows_search_disable"
	xboxServicesDisableBooster "github.com/oLenador/mulltbost/internal/boosters/fpsboost/boosters/xbox_services_disable"

	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func GetAllPlugins() []inbound.BoosterUseCase {
	return []inbound.BoosterUseCase{
		appBackgroundDisableBooster.NewAppBackgroundDisable()(),
		bluetoothDisableBooster.NewBluetoothDisableBooster(),
		cpuAffinityBooster.NewCpuAffinityBooster(),
		cpuGamePriorityBooster.NewCpuGamePriorityBooster(),
		cpuMitigationsBooster.NewCpuMitigationsBooster(),
		cpuSystemServicesReduceBooster.NewCpuSystemServicesReduceBooster(),
		explorerOptimizeBooster.NewExplorerOptimizeBooster(),
		errorReportingDisableBooster.NewErrorReportingDisableBooster(),
		fastStartupDisableBooster.NewFastStartupDisableBooster(),
		gameDvrDisableBooster.NewGameDvrDisableBooster(),
		gpuLowLatencyModeBooster.NewGpuLowLatencyModeBooster(),
		gpuSchedulingBooster.NewGpuSchedulingBooster(),
		gpuSchedulingHardwareBooster.NewGpuSchedulingHardwareBooster(),
		gpuIrqOptimizeBooster.NewGpuIrqOptimizeBooster(),
		gpuPowerPolicyBooster.NewGpuPowerPolicyBooster(),
		hibernationDisableBooster.NewHibernationDisableBooster(),
		mapsManagerDisableBooster.NewMapsManagerDisableBooster(),
		memoryAllocationBooster.NewMemoryAllocationBooster(),
		memoryCacheOptimizeBooster.NewMemoryCacheOptimizeBooster(),
		memoryFlushBooster.NewMemoryFlushBooster(),
		memoryFreeBooster.NewMemoryFreeBooster(),
		multBoostPowerBooster.NewMultBoostPowerBooster(),
		nonPagedPoolOptimizeBooster.NewNonPagedPoolOptimizeBooster(),
		pageFileSizeBooster.NewPageFileSizeBooster(),
		prefetchSuperfetchDisableBooster.NewPrefetchSuperfetchDisableBooster(),
		processShutdownDelayBooster.NewProcessShutdownDelayBooster(),
		searchSuggestionsDisableBooster.NewSearchSuggestionsDisableBooster(),
		servicesDisableBooster.NewServicesDisableBooster(),
		smartScreenDisableBooster.NewSmartScreenDisableBooster(),
		standbyListCleanBooster.NewStandbyListCleanBooster(),
		sysmainOptimizeBooster.NewSysmainOptimizeBooster(),
		telemetryDisableBooster.NewTelemetryDisableBooster(),
		threadSchedulingBooster.NewThreadSchedulingBooster(),
		timerResolutionBooster.NewTimerResolutionBooster(),
		trimSsdBooster.NewTrimSsdBooster(),
		visualEffectsDisableBooster.NewVisualEffectsDisableBooster(),
		windowsDefenderDisableBooster.NewWindowsDefenderDisableBooster(),
		windowsSearchDisableBooster.NewWindowsSearchDisableBooster(),
		xboxServicesDisableBooster.NewXboxServicesDisableBooster(),
	}
}

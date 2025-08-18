package precision

import (
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
	audioLatencyBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/precision/boosters/audio_latency"
	controllerPrecisionBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/precision/boosters/controller_precision"
	displayPrecisionBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/precision/boosters/display_precision"
	keyboardPrecisionBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/precision/boosters/keyboard_precision"
	mousePrecisionBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/precision/boosters/mouse_precision"
	powerSystemUsbBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/precision/boosters/power_system_usb"
)

func GetAllPlugins() []inbound.BoosterUseCase {
	return []inbound.BoosterUseCase{
		powerSystemUsbBooster.NewPowerSystemUSBBooster(),
		audioLatencyBooster.NewAudioLatencyBooster(),
		controllerPrecisionBooster.NewControllerPrecisionBooster(),
		displayPrecisionBooster.NewDisplayPrecisionBooster(),
		keyboardPrecisionBooster.NewKeyboardPrecisionBooster(),
		mousePrecisionBooster.NewMousePrecisionBooster(),
	}
}

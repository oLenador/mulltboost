package precision

import (
	powerSystemUsbBooster "github.com/oLenador/mulltbost/internal/boosters/precision/boosters/power_system_usb"
	audioLatencyBooster "github.com/oLenador/mulltbost/internal/boosters/precision/boosters/audio_latency"
	controllerPrecisionBooster "github.com/oLenador/mulltbost/internal/boosters/precision/boosters/controller_precision"
	displayPrecisionBooster "github.com/oLenador/mulltbost/internal/boosters/precision/boosters/display_precision"
	keyboardPrecisionBooster "github.com/oLenador/mulltbost/internal/boosters/precision/boosters/keyboard_precision"
	mousePrecisionBooster "github.com/oLenador/mulltbost/internal/boosters/precision/boosters/mouse_precision"

	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
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

package outbound

import (
	"context"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type InputDeviceService interface {
	// Dispositivos de entrada
	GetInputDevices(ctx context.Context) ([]*entities.InputDevice, error)
	
	// Status
	GetInputDeviceConfiguration(ctx context.Context, deviceID string) (*entities.InputDevice, error)
}

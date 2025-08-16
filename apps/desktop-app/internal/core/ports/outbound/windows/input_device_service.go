package windows

import (
	"context"
)

type InputDeviceService interface {
	// Dispositivos de entrada
	GetInputDevices(ctx context.Context) ([]*entities.InputDevice, error)
	OptimizeMouseSettings(ctx context.Context) error
	OptimizeKeyboardSettings(ctx context.Context) error
	
	// Configurações
	SetMouseSensitivity(ctx context.Context, sensitivity float64) error
	SetKeyboardRepeatRate(ctx context.Context, rate int) error
	DisableMouseAcceleration(ctx context.Context) error
	
	// Status
	IsInputOptimized(ctx context.Context) (bool, error)
	GetInputDeviceConfiguration(ctx context.Context, deviceID string) (*entities.InputDevice, error)
}

package entities

import "time"

// DeviceType representa os tipos de dispositivo de entrada
type DeviceType string

const (
	DeviceTypeMouse    DeviceType = "mouse"
	DeviceTypeKeyboard DeviceType = "keyboard"
	DeviceTypeOther    DeviceType = "other"
	DeviceTypeGamepad  DeviceType = "gamepad"
	DeviceTypeTablet   DeviceType = "tablet"
	DeviceTypeTouchpad DeviceType = "touchpad"
)

// InputDevice representa um dispositivo de entrada do sistema
type InputDevice struct {
	ID           string                 `json:"id" bson:"_id,omitempty"`
	Name         string                 `json:"name" bson:"name"`
	DeviceType   DeviceType             `json:"device_type" bson:"device_type"`
	Description  string                 `json:"description" bson:"description"`
	Manufacturer string                 `json:"manufacturer" bson:"manufacturer"`
	HardwareID   string                 `json:"hardware_id" bson:"hardware_id"`
	IsConnected  bool                   `json:"is_connected" bson:"is_connected"`
	Status       string                 `json:"status" bson:"status"`
	Properties   map[string]interface{} `json:"properties" bson:"properties"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewInputDevice cria uma nova instância de InputDevice
func NewInputDevice(name string, deviceType DeviceType) *InputDevice {
	now := time.Now()
	return &InputDevice{
		Name:        name,
		DeviceType:  deviceType,
		IsConnected: false,
		Properties:  make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsActive verifica se o dispositivo está ativo e conectado
func (d *InputDevice) IsActive() bool {
	return d.IsConnected && d.Status == "OK"
}

// GetProperty retorna uma propriedade específica do dispositivo
func (d *InputDevice) GetProperty(key string) interface{} {
	if d.Properties == nil {
		return nil
	}
	return d.Properties[key]
}

// SetProperty define uma propriedade do dispositivo
func (d *InputDevice) SetProperty(key string, value interface{}) {
	if d.Properties == nil {
		d.Properties = make(map[string]interface{})
	}
	d.Properties[key] = value
	d.UpdatedAt = time.Now()
}

// UpdateStatus atualiza o status do dispositivo
func (d *InputDevice) UpdateStatus(status string, connected bool) {
	d.Status = status
	d.IsConnected = connected
	d.UpdatedAt = time.Now()
}

// String implementa a interface Stringer
func (d *InputDevice) String() string {
	return d.Name
}
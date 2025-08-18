package system

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type WindowsInputDeviceService struct {
	comInitialized bool
}

func NewWindowsInputDeviceService() *WindowsInputDeviceService {
	return &WindowsInputDeviceService{}
}

func (s *WindowsInputDeviceService) initializeCOM() error {
	if !s.comInitialized {
		err := ole.CoInitialize(0)
		if err != nil {
			return fmt.Errorf("failed to initialize COM: %w", err)
		}
		s.comInitialized = true
	}
	return nil
}

// Cleanup limpa recursos COM
func (s *WindowsInputDeviceService) Cleanup() {
	if s.comInitialized {
		ole.CoUninitialize()
		s.comInitialized = false
	}
}

// GetInputDevices retorna todos os dispositivos de entrada detectados
func (s *WindowsInputDeviceService) GetInputDevices(ctx context.Context) ([]*entities.InputDevice, error) {
	if err := s.initializeCOM(); err != nil {
		return nil, err
	}

	var devices []*entities.InputDevice

	// Buscar dispositivos de mouse
	mouseDevices, err := s.getMouseDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting mouse devices: %w", err)
	}
	devices = append(devices, mouseDevices...)

	// Buscar dispositivos de teclado
	keyboardDevices, err := s.getKeyboardDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting keyboard devices: %w", err)
	}
	devices = append(devices, keyboardDevices...)

	// Buscar outros dispositivos HID
	hidDevices, err := s.getHIDDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting HID devices: %w", err)
	}
	devices = append(devices, hidDevices...)

	return devices, nil
}

// GetInputDeviceConfiguration retorna a configuração de um dispositivo específico
func (s *WindowsInputDeviceService) GetInputDeviceConfiguration(ctx context.Context, deviceID string) (*entities.InputDevice, error) {
	if err := s.initializeCOM(); err != nil {
		return nil, err
	}

	devices, err := s.GetInputDevices(ctx)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.ID == deviceID {
			return device, nil
		}
	}

	return nil, fmt.Errorf("device with ID %s not found", deviceID)
}

// getMouseDevices busca dispositivos de mouse via WMI
func (s *WindowsInputDeviceService) getMouseDevices(ctx context.Context) ([]*entities.InputDevice, error) {
	query := "SELECT * FROM Win32_PointingDevice"
	results, err := s.executeWMIQuery(query)
	if err != nil {
		return nil, err
	}
	defer results.Release()

	var devices []*entities.InputDevice
	
	// Iterar através dos resultados usando oleutil
	err = oleutil.ForEach(results, func(v *ole.VARIANT) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if v.VT != ole.VT_DISPATCH {
			return nil
		}

		dispatch := v.ToIDispatch()
		defer dispatch.Release()

		device := s.parseMouseDevice(dispatch)
		if device != nil {
			devices = append(devices, device)
		}
		
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating mouse devices: %w", err)
	}

	return devices, nil
}

// getKeyboardDevices busca dispositivos de teclado via WMI
func (s *WindowsInputDeviceService) getKeyboardDevices(ctx context.Context) ([]*entities.InputDevice, error) {
	query := "SELECT * FROM Win32_Keyboard"
	results, err := s.executeWMIQuery(query)
	if err != nil {
		return nil, err
	}
	defer results.Release()

	var devices []*entities.InputDevice
	
	// Iterar através dos resultados usando oleutil
	err = oleutil.ForEach(results, func(v *ole.VARIANT) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if v.VT != ole.VT_DISPATCH {
			return nil
		}

		dispatch := v.ToIDispatch()
		defer dispatch.Release()

		device := s.parseKeyboardDevice(dispatch)
		if device != nil {
			devices = append(devices, device)
		}
		
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating keyboard devices: %w", err)
	}

	return devices, nil
}

// getHIDDevices busca outros dispositivos HID via WMI
func (s *WindowsInputDeviceService) getHIDDevices(ctx context.Context) ([]*entities.InputDevice, error) {
	// Query mais específica para dispositivos HID reais
	query := "SELECT * FROM Win32_PnPEntity WHERE PNPClass = 'HIDClass' OR Service = 'HidUsb'"
	results, err := s.executeWMIQuery(query)
	if err != nil {
		return nil, err
	}
	defer results.Release()

	var devices []*entities.InputDevice
	
	// Iterar através dos resultados usando oleutil
	err = oleutil.ForEach(results, func(v *ole.VARIANT) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if v.VT != ole.VT_DISPATCH {
			return nil
		}

		dispatch := v.ToIDispatch()
		defer dispatch.Release()

		device := s.parseHIDDevice(dispatch)
		if device != nil {
			devices = append(devices, device)
		}
		
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating HID devices: %w", err)
	}

	return devices, nil
}

// executeWMIQuery executa uma consulta WMI e retorna os resultados
func (s *WindowsInputDeviceService) executeWMIQuery(query string) (*ole.IDispatch, error) {
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, fmt.Errorf("failed to create WMI locator: %w", err)
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to get WMI dispatch: %w", err)
	}
	defer wmi.Release()

	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WMI service: %w", err)
	}
	defer serviceRaw.Clear()

	service := serviceRaw.ToIDispatch()
	defer service.Release()

	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute WMI query: %w", err)
	}
	defer resultRaw.Clear()

	// Retorna a interface IDispatch diretamente
	return resultRaw.ToIDispatch(), nil
}

// parseMouseDevice converte um objeto WMI Win32_PointingDevice em InputDevice
func (s *WindowsInputDeviceService) parseMouseDevice(dispatch *ole.IDispatch) *entities.InputDevice {
	device := &entities.InputDevice{
		DeviceType: "mouse",
		Properties: make(map[string]interface{}),
	}

	// Extrair informações básicas
	if name := s.getStringProperty(dispatch, "Name"); name != "" {
		device.Name = name
		device.ID = s.generateDeviceID("mouse", name)
	} else {
		return nil // Dispositivo inválido sem nome
	}

	if description := s.getStringProperty(dispatch, "Description"); description != "" {
		device.Description = description
	}

	if manufacturer := s.getStringProperty(dispatch, "Manufacturer"); manufacturer != "" {
		device.Manufacturer = manufacturer
	}

	if pnpDeviceID := s.getStringProperty(dispatch, "PNPDeviceID"); pnpDeviceID != "" {
		device.HardwareID = pnpDeviceID
	}

	// Status do dispositivo
	status := s.getStringProperty(dispatch, "Status")
	if status != "" {
		device.IsConnected = strings.EqualFold(status, "OK")
		device.Status = status
	} else {
		// Assumir conectado se não especificado
		device.IsConnected = true
		device.Status = "OK"
	}

	// Configurações específicas do mouse
	if handedness := s.getStringProperty(dispatch, "Handedness"); handedness != "" {
		device.Properties["handedness"] = handedness
	}
	if numberOfButtons := s.getUint32Property(dispatch, "NumberOfButtons"); numberOfButtons > 0 {
		device.Properties["number_of_buttons"] = numberOfButtons
	}
	if pointingType := s.getStringProperty(dispatch, "PointingType"); pointingType != "" {
		device.Properties["pointing_type"] = pointingType
	}
	if resolution := s.getUint32Property(dispatch, "Resolution"); resolution > 0 {
		device.Properties["resolution"] = resolution
	}

	return device
}

// parseKeyboardDevice converte um objeto WMI Win32_Keyboard em InputDevice
func (s *WindowsInputDeviceService) parseKeyboardDevice(dispatch *ole.IDispatch) *entities.InputDevice {
	device := &entities.InputDevice{
		DeviceType: "keyboard",
		Properties: make(map[string]interface{}),
	}

	// Extrair informações básicas
	if name := s.getStringProperty(dispatch, "Name"); name != "" {
		device.Name = name
		device.ID = s.generateDeviceID("keyboard", name)
	} else {
		return nil // Dispositivo inválido sem nome
	}

	if description := s.getStringProperty(dispatch, "Description"); description != "" {
		device.Description = description
	}

	if pnpDeviceID := s.getStringProperty(dispatch, "PNPDeviceID"); pnpDeviceID != "" {
		device.HardwareID = pnpDeviceID
	}

	// Status do dispositivo
	status := s.getStringProperty(dispatch, "Status")
	if status != "" {
		device.IsConnected = strings.EqualFold(status, "OK")
		device.Status = status
	} else {
		// Assumir conectado se não especificado
		device.IsConnected = true
		device.Status = "OK"
	}

	// Configurações específicas do teclado
	if layout := s.getStringProperty(dispatch, "Layout"); layout != "" {
		device.Properties["layout"] = layout
	}
	if numberOfFunctionKeys := s.getUint32Property(dispatch, "NumberOfFunctionKeys"); numberOfFunctionKeys > 0 {
		device.Properties["number_of_function_keys"] = numberOfFunctionKeys
	}
	if password := s.getBoolProperty(dispatch, "Password"); password {
		device.Properties["password_enabled"] = password
	}

	return device
}

// parseHIDDevice converte um objeto WMI Win32_PnPEntity em InputDevice (para dispositivos HID)
func (s *WindowsInputDeviceService) parseHIDDevice(dispatch *ole.IDispatch) *entities.InputDevice {
	name := s.getStringProperty(dispatch, "Name")
	if name == "" || s.isSystemDevice(name) {
		return nil // Pular dispositivos de sistema que não são de entrada
	}

	device := &entities.InputDevice{
		DeviceType: "other",
		Properties: make(map[string]interface{}),
	}

	device.Name = name
	device.ID = s.generateDeviceID("hid", name)

	if description := s.getStringProperty(dispatch, "Description"); description != "" {
		device.Description = description
	}

	if manufacturer := s.getStringProperty(dispatch, "Manufacturer"); manufacturer != "" {
		device.Manufacturer = manufacturer
	}

	if pnpDeviceID := s.getStringProperty(dispatch, "PNPDeviceID"); pnpDeviceID != "" {
		device.HardwareID = pnpDeviceID
	}

	// Status do dispositivo
	status := s.getStringProperty(dispatch, "Status")
	if status != "" {
		device.IsConnected = strings.EqualFold(status, "OK")
		device.Status = status
	} else {
		// Assumir conectado se não especificado
		device.IsConnected = true
		device.Status = "OK"
	}

	// Propriedades específicas do HID
	if service := s.getStringProperty(dispatch, "Service"); service != "" {
		device.Properties["service"] = service
	}
	if pnpClass := s.getStringProperty(dispatch, "PNPClass"); pnpClass != "" {
		device.Properties["pnp_class"] = pnpClass
	}

	return device
}

// Métodos auxiliares para extrair propriedades com tratamento de erro melhorado
func (s *WindowsInputDeviceService) getStringProperty(dispatch *ole.IDispatch, property string) string {
	result, err := oleutil.GetProperty(dispatch, property)
	if err != nil {
		return ""
	}
	defer result.Clear()

	if result.VT == ole.VT_NULL || result.VT == ole.VT_EMPTY {
		return ""
	}

	return result.ToString()
}

func (s *WindowsInputDeviceService) getUint32Property(dispatch *ole.IDispatch, property string) uint32 {
	result, err := oleutil.GetProperty(dispatch, property)
	if err != nil {
		return 0
	}
	defer result.Clear()

	if result.VT == ole.VT_NULL || result.VT == ole.VT_EMPTY {
		return 0
	}

	switch result.VT {
	case ole.VT_I4, ole.VT_UI4:
		return uint32(result.Val)
	case ole.VT_I2, ole.VT_UI2:
		return uint32(result.Val)
	case ole.VT_BSTR:
		val, err := strconv.ParseUint(result.ToString(), 10, 32)
		if err != nil {
			return 0
		}
		return uint32(val)
	default:
		return 0
	}
}

func (s *WindowsInputDeviceService) getBoolProperty(dispatch *ole.IDispatch, property string) bool {
	result, err := oleutil.GetProperty(dispatch, property)
	if err != nil {
		return false
	}
	defer result.Clear()

	if result.VT == ole.VT_NULL || result.VT == ole.VT_EMPTY {
		return false
	}

	switch result.VT {
	case ole.VT_BOOL:
		return result.Val != 0
	case ole.VT_BSTR:
		return strings.EqualFold(result.ToString(), "true")
	default:
		return false
	}
}

// generateDeviceID gera um ID único para o dispositivo
func (s *WindowsInputDeviceService) generateDeviceID(deviceType, name string) string {
	// Remove caracteres especiais e espaços
	cleanName := strings.ReplaceAll(strings.ToLower(name), " ", "_")
	cleanName = strings.ReplaceAll(cleanName, ".", "")
	cleanName = strings.ReplaceAll(cleanName, "-", "_")
	cleanName = strings.ReplaceAll(cleanName, "(", "")
	cleanName = strings.ReplaceAll(cleanName, ")", "")
	cleanName = strings.ReplaceAll(cleanName, "&", "and")
	
	// Remove múltiplos underscores consecutivos
	for strings.Contains(cleanName, "__") {
		cleanName = strings.ReplaceAll(cleanName, "__", "_")
	}
	
	// Remove underscores no início e fim
	cleanName = strings.Trim(cleanName, "_")
	
	return fmt.Sprintf("%s_%s", deviceType, cleanName)
}

// isSystemDevice verifica se é um dispositivo de sistema que deve ser ignorado
func (s *WindowsInputDeviceService) isSystemDevice(name string) bool {
	systemDevices := []string{
		"microsoft ps/2 mouse",
		"hid-compliant mouse",
		"standard ps/2 keyboard",
		"hid keyboard device",
		"system speaker",
		"beep",
		"microsoft system management bios",
		"acpi",
		"root hub",
		"composite device",
	}

	lowerName := strings.ToLower(name)
	for _, sysDevice := range systemDevices {
		if strings.Contains(lowerName, sysDevice) {
			return false // Na verdade, estes SÃO dispositivos de entrada válidos
		}
	}

	// Filtrar apenas drivers realmente não relacionados
	systemDrivers := []string{
		"audio", "video", "display", "network", "storage",
		"disk", "cdrom", "printer", "modem", "bluetooth radio",
	}

	for _, sysDriver := range systemDrivers {
		if strings.Contains(lowerName, sysDriver) {
			return true
		}
	}

	return false
}

// Implementação para uso como singleton ou service locator
var defaultService *WindowsInputDeviceService

// GetDefaultService retorna a instância padrão do serviço
func GetDefaultService() *WindowsInputDeviceService {
	if defaultService == nil {
		defaultService = NewWindowsInputDeviceService()
	}
	return defaultService
}

// CleanupDefaultService limpa a instância padrão
func CleanupDefaultService() {
	if defaultService != nil {
		defaultService.Cleanup()
		defaultService = nil
	}
}
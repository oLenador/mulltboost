//go:build windows
// +build windows

// ===== ARQUIVO: internal/adapters/system/windows_network_adapter_service.go =====
package system

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"golang.org/x/sys/windows/registry"
)

const (
	// Timeouts de segurança
	wmiQueryTimeout     = 30 * time.Second
	commandTimeout      = 60 * time.Second
	maxConcurrentOps    = 5
	
	// Limites de validação
	minPriority         = 1
	maxPriority         = 9999
	maxAdapterNameLen   = 256
	maxDescriptionLen   = 512
)

// WindowsNetworkAdapterService implementa NetworkAdapterService para Windows
type WindowsNetworkAdapterService struct {
	ole       *ole.IUnknown
	mu        sync.RWMutex
	semaphore chan struct{} // Limita operações concorrentes
	closed    bool
}

// NewWindowsNetworkAdapterService cria uma nova instância do serviço
func NewWindowsNetworkAdapterService() (*WindowsNetworkAdapterService, error) {
	if err := ole.CoInitialize(0); err != nil {
		return nil, fmt.Errorf("failed to initialize OLE: %w", err)
	}

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		ole.CoUninitialize()
		return nil, fmt.Errorf("failed to create WMI locator: %w", err)
	}

	return &WindowsNetworkAdapterService{
		ole:       unknown,
		semaphore: make(chan struct{}, maxConcurrentOps),
	}, nil
}

// Close libera recursos de forma segura
func (s *WindowsNetworkAdapterService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.closed {
		return nil
	}
	
	s.closed = true
	
	if s.ole != nil {
		s.ole.Release()
		s.ole = nil
	}
	
	ole.CoUninitialize()
	close(s.semaphore)
	
	return nil
}

// GetNetworkAdapters retorna todos os adaptadores de rede do sistema
func (s *WindowsNetworkAdapterService) GetNetworkAdapters(ctx context.Context) ([]*entities.NetworkAdapter, error) {
	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Timeout para operação WMI
	ctx, cancel := context.WithTimeout(ctx, wmiQueryTimeout)
	defer cancel()

	service, err := s.connectWMI()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WMI: %w", err)
	}
	defer func() {
		if service != nil {
			service.Release()
		}
	}()

	// Query segura para adaptadores de rede físicos e virtuais ativos
	query := "SELECT Name, Description, MACAddress, DeviceID, AdapterType, NetConnectionStatus, Speed FROM Win32_NetworkAdapter WHERE NetEnabled = TRUE AND MACAddress IS NOT NULL"
	
	result, err := s.executeWMIQuery(ctx, service, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute WMI query: %w", err)
	}
	defer result.Clear()

	collection := result.ToIDispatch()
	defer collection.Release()

	// Conta os items com verificação de erro
	countVar, err := oleutil.GetProperty(collection, "Count")
	if err != nil {
		return nil, fmt.Errorf("failed to get collection count: %w", err)
	}
	count := int(countVar.Val)

	if count < 0 || count > 1000 { // Limite de segurança
		return nil, fmt.Errorf("invalid adapter count: %d", count)
	}

	var adapters []*entities.NetworkAdapter
	
	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			return adapters, ctx.Err() // Retorna o que foi coletado até agora
		default:
		}

		itemVar, err := oleutil.CallMethod(collection, "ItemIndex", i)
		if err != nil {
			continue // Log error em produção, mas continua
		}

		item := itemVar.ToIDispatch()
		adapter, err := s.parseNetworkAdapter(item)
		item.Release()

		if err != nil {
			continue // Log error em produção, mas continua
		}

		if adapter != nil && adapter.Validate() == nil {
			adapters = append(adapters, adapter)
		}
	}

	return adapters, nil
}

// parseNetworkAdapter converte um objeto WMI em NetworkAdapter com validação
func (s *WindowsNetworkAdapterService) parseNetworkAdapter(item *ole.IDispatch) (*entities.NetworkAdapter, error) {
	if item == nil {
		return nil, fmt.Errorf("invalid WMI item")
	}

	// Propriedades básicas com sanitização
	name, _ := s.getStringProperty(item, "Name")
	name = s.sanitizeString(name, maxAdapterNameLen)
	
	description, _ := s.getStringProperty(item, "Description")
	description = s.sanitizeString(description, maxDescriptionLen)
	
	macAddress, _ := s.getStringProperty(item, "MACAddress")
	macAddress = s.sanitizeMACAddress(macAddress)
	
	deviceID, _ := s.getStringProperty(item, "DeviceID")
	deviceID = s.sanitizeString(deviceID, 64)
	
	adapterType, _ := s.getStringProperty(item, "AdapterType")
	
	// Status e velocidade
	netConnectionStatus, _ := s.getUint32Property(item, "NetConnectionStatus")
	speed, _ := s.getUint64Property(item, "Speed")

	// Validações básicas
	if name == "" || deviceID == "" {
		return nil, fmt.Errorf("invalid adapter data: missing required fields")
	}

	// Determina o tipo do adaptador
	adapterTypeEnum := s.determineAdapterType(adapterType, description)

	// Status de conexão
	status := s.parseConnectionStatus(netConnectionStatus)

	adapter := &entities.NetworkAdapter{
		ID:          deviceID,
		Name:        name,
		Description: description,
		MACAddress:  macAddress,
		Type:        adapterTypeEnum,
		Status:      status,
		Speed:       speed,
		IsEnabled:   status == entities.NetworkStatusConnected,
		Metrics: &entities.NetworkMetrics{
			BytesSent:       0,
			BytesReceived:   0,
			PacketsSent:     0,
			PacketsReceived: 0,
			Errors:          0,
			Dropped:         0,
			LastUpdated:     time.Now(),
		},
	}

	return adapter, nil
}

// GetNetworkStatistics retorna estatísticas detalhadas com timeout
func (s *WindowsNetworkAdapterService) GetNetworkStatistics(ctx context.Context, adapterID string) (*entities.NetworkAdapter, error) {
	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	if err := s.validateAdapterID(adapterID); err != nil {
		return nil, fmt.Errorf("invalid adapter ID: %w", err)
	}

	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	ctx, cancel := context.WithTimeout(ctx, wmiQueryTimeout)
	defer cancel()

	service, err := s.connectWMI()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WMI: %w", err)
	}
	defer service.Release()

	// Query segura com escape de caracteres especiais
	escapedID := s.escapeWMIString(adapterID)
	query := fmt.Sprintf("SELECT * FROM Win32_PerfRawData_Tcpip_NetworkInterface WHERE Name LIKE '%%%s%%'", escapedID)
	
	result, err := s.executeWMIQuery(ctx, service, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statistics query: %w", err)
	}
	defer result.Clear()

	collection := result.ToIDispatch()
	defer collection.Release()

	countVar, _ := oleutil.GetProperty(collection, "Count")
	if int(countVar.Val) == 0 {
		return nil, fmt.Errorf("adapter statistics not found: %s", adapterID)
	}

	itemVar, _ := oleutil.CallMethod(collection, "ItemIndex", 0)
	item := itemVar.ToIDispatch()
	defer item.Release()

	// Coleta estatísticas
	bytesSent, _ := s.getUint64Property(item, "BytesSentPerSec")
	bytesReceived, _ := s.getUint64Property(item, "BytesReceivedPerSec")
	packetsSent, _ := s.getUint64Property(item, "PacketsOutboundErrors")
	packetsReceived, _ := s.getUint64Property(item, "PacketsReceivedErrors")
	errors, _ := s.getUint64Property(item, "PacketsReceivedErrors")

	// Busca informações básicas do adaptador
	adapter, err := s.getBasicAdapterInfo(ctx, adapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic adapter info: %w", err)
	}

	// Atualiza métricas
	adapter.Metrics = &entities.NetworkMetrics{
		BytesSent:       bytesSent,
		BytesReceived:   bytesReceived,
		PacketsSent:     packetsSent,
		PacketsReceived: packetsReceived,
		Errors:          errors,
		Dropped:         0,
		LastUpdated:     time.Now(),
	}

	return adapter, nil
}

// ResetNetworkStack executa comandos com timeout e verificação de privilégios
func (s *WindowsNetworkAdapterService) ResetNetworkStack(ctx context.Context) error {
	if err := s.checkClosed(); err != nil {
		return err
	}

	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return ctx.Err()
	}

	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	commands := []struct {
		name string
		args []string
	}{
		{"netsh", []string{"winsock", "reset"}},
		{"netsh", []string{"int", "ip", "reset"}},
		{"ipconfig", []string{"/flushdns"}},
		{"netsh", []string{"int", "tcp", "reset"}},
	}

	for _, cmd := range commands {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := s.runElevatedCommand(ctx, cmd.name, cmd.args...); err != nil {
			return fmt.Errorf("failed to execute %s %v: %w", cmd.name, cmd.args, err)
		}
	}

	return nil
}

// SetNetworkPriority com validação aprimorada
func (s *WindowsNetworkAdapterService) SetNetworkPriority(ctx context.Context, adapterID string, priority int) error {
	if err := s.checkClosed(); err != nil {
		return err
	}

	if err := s.validateAdapterID(adapterID); err != nil {
		return fmt.Errorf("invalid adapter ID: %w", err)
	}

	if priority < minPriority || priority > maxPriority {
		return fmt.Errorf("invalid priority: %d (must be between %d and %d)", priority, minPriority, maxPriority)
	}

	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return ctx.Err()
	}

	// Sanitiza o adapterID para prevenir path traversal
	cleanAdapterID := s.sanitizeRegistryPath(adapterID)
	keyPath := fmt.Sprintf(`SYSTEM\CurrentControlSet\Control\Class\{4d36e972-e325-11ce-bfc1-08002be10318}\%s`, cleanAdapterID)
	
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open adapter registry key (check permissions): %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("NetworkInterfaceMetric", uint32(priority)); err != nil {
		return fmt.Errorf("failed to set network priority: %w", err)
	}

	return nil
}

// DisableNetworkThrottling com verificação de permissões
func (s *WindowsNetworkAdapterService) DisableNetworkThrottling(ctx context.Context) error {
	if err := s.checkClosed(); err != nil {
		return err
	}

	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }()
	case <-ctx.Done():
		return ctx.Err()
	}

	// Configurações validadas
	settings := map[string]uint32{
		"NetworkThrottlingIndex": 0xffffffff,
		"SystemResponsiveness":   0,
		"NetworkInterfaceMetric": 1,
	}

	keyPath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Multimedia\SystemProfile\Tasks\Games`
	if err := s.applyRegistrySettings(keyPath, settings); err != nil {
		return fmt.Errorf("failed to disable throttling: %w", err)
	}

	// Configurações TCP/IP globais
	globalSettings := map[string]uint32{
		"TCPNoDelay":                    1,
		"TcpAckFrequency":              1,
		"TcpDelAckTicks":               0,
		"TCPInitialRTT":                300,
		"TcpMaxDataRetransmissions":    3,
	}

	globalKeyPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`
	if err := s.applyRegistrySettings(globalKeyPath, globalSettings); err != nil {
		// Log warning mas não falha para configurações opcionais
		return nil
	}

	return nil
}

// IsNetworkOptimized com verificações de segurança
func (s *WindowsNetworkAdapterService) IsNetworkOptimized(ctx context.Context) (bool, error) {
	if err := s.checkClosed(); err != nil {
		return false, err
	}

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	checks := []func() bool{
		s.isThrottlingDisabled,
		s.isHighPrioritySet,
		s.isTCPOptimized,
	}

	for _, check := range checks {
		if !check() {
			return false, nil
		}
	}

	return true, nil
}

// ===== MÉTODOS AUXILIARES DE SEGURANÇA =====

func (s *WindowsNetworkAdapterService) checkClosed() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return fmt.Errorf("service is closed")
	}
	return nil
}

func (s *WindowsNetworkAdapterService) connectWMI() (*ole.IDispatch, error) {
	if s.ole == nil {
		return nil, fmt.Errorf("OLE not initialized")
	}

	locator, err := oleutil.QueryInterface(s.ole, ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to query interface: %w", err)
	}

	serviceRaw, err := oleutil.CallMethod(locator, "ConnectServer", nil, `root\cimv2`)
	if err != nil {
		locator.Release()
		return nil, fmt.Errorf("failed to connect to WMI service: %w", err)
	}

	service := serviceRaw.ToIDispatch()
	locator.Release()

	return service, nil
}

func (s *WindowsNetworkAdapterService) executeWMIQuery(ctx context.Context, service *ole.IDispatch, query string) (*ole.VARIANT, error) {
	done := make(chan struct{})
	var result *ole.VARIANT
	var err error

	go func() {
		defer close(done)
		resultVar, queryErr := oleutil.CallMethod(service, "ExecQuery", query)
		result = &resultVar
		err = queryErr
	}()

	select {
	case <-done:
		return result, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *WindowsNetworkAdapterService) sanitizeString(input string, maxLen int) string {
	if len(input) > maxLen {
		input = input[:maxLen]
	}
	// Remove caracteres de controle
	return strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, input)
}

func (s *WindowsNetworkAdapterService) sanitizeMACAddress(mac string) string {
	// Remove tudo exceto dígitos hexadecimais e separadores comuns
	cleaned := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f') || r == ':' || r == '-' {
			return r
		}
		return -1
	}, mac)
	
	if len(cleaned) > 17 { // Máximo para formato XX:XX:XX:XX:XX:XX
		cleaned = cleaned[:17]
	}
	
	return cleaned
}

func (s *WindowsNetworkAdapterService) validateAdapterID(adapterID string) error {
	if adapterID == "" {
		return fmt.Errorf("adapter ID cannot be empty")
	}
	
	if len(adapterID) > 64 {
		return fmt.Errorf("adapter ID too long")
	}
	
	// Verifica caracteres perigosos
	dangerous := []string{"../", "..\\", "|", "&", ";", "$", "`", "\"", "'"}
	for _, d := range dangerous {
		if strings.Contains(adapterID, d) {
			return fmt.Errorf("adapter ID contains invalid characters")
		}
	}
	
	return nil
}

func (s *WindowsNetworkAdapterService) sanitizeRegistryPath(path string) string {
	// Remove caracteres perigosos para registry paths
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return -1
	}, path)
	
	if len(safe) > 32 {
		safe = safe[:32]
	}
	
	return safe
}

func (s *WindowsNetworkAdapterService) escapeWMIString(input string) string {
	// Escapa caracteres especiais para queries WMI
	input = strings.ReplaceAll(input, "'", "''")
	input = strings.ReplaceAll(input, "\\", "\\\\")
	input = strings.ReplaceAll(input, "%", "\\%")
	input = strings.ReplaceAll(input, "_", "\\_")
	return input
}

func (s *WindowsNetworkAdapterService) applyRegistrySettings(keyPath string, settings map[string]uint32) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.SET_VALUE)
	if err != nil {
		// Tenta criar se não existir
		key, _, err = registry.CreateKey(registry.LOCAL_MACHINE, keyPath, registry.SET_VALUE)
		if err != nil {
			return fmt.Errorf("failed to open/create registry key %s: %w", keyPath, err)
		}
	}
	defer key.Close()

	for name, value := range settings {
		if err := key.SetDWordValue(name, value); err != nil {
			return fmt.Errorf("failed to set registry value %s: %w", name, err)
		}
	}

	return nil
}

// ===== MÉTODOS AUXILIARES ORIGINAIS (mantidos por compatibilidade) =====

func (s *WindowsNetworkAdapterService) getStringProperty(item *ole.IDispatch, property string) (string, error) {
	if item == nil {
		return "", fmt.Errorf("invalid dispatch item")
	}
	
	prop, err := oleutil.GetProperty(item, property)
	if err != nil {
		return "", err
	}
	
	if prop.VT == ole.VT_NULL || prop.VT == ole.VT_EMPTY {
		return "", nil
	}
	
	return prop.ToString(), nil
}

func (s *WindowsNetworkAdapterService) getUint32Property(item *ole.IDispatch, property string) (uint32, error) {
	if item == nil {
		return 0, fmt.Errorf("invalid dispatch item")
	}
	
	prop, err := oleutil.GetProperty(item, property)
	if err != nil {
		return 0, err
	}
	
	if prop.VT == ole.VT_NULL || prop.VT == ole.VT_EMPTY {
		return 0, nil
	}
	
	return uint32(prop.Val), nil
}

func (s *WindowsNetworkAdapterService) getUint64Property(item *ole.IDispatch, property string) (uint64, error) {
	if item == nil {
		return 0, fmt.Errorf("invalid dispatch item")
	}
	
	prop, err := oleutil.GetProperty(item, property)
	if err != nil {
		return 0, err
	}
	
	if prop.VT == ole.VT_NULL || prop.VT == ole.VT_EMPTY {
		return 0, nil
	}
	
	// Handling diferentes tipos de VARIANT
	switch prop.VT {
	case ole.VT_UI8:
		return uint64(prop.Val), nil
	case ole.VT_I8:
		return uint64(prop.Val), nil
	case ole.VT_UI4:
		return uint64(uint32(prop.Val)), nil
	case ole.VT_I4:
		return uint64(int32(prop.Val)), nil
	default:
		return uint64(prop.Val), nil
	}
}

func (s *WindowsNetworkAdapterService) determineAdapterType(adapterType, description string) entities.NetworkAdapterType {
	description = strings.ToLower(strings.TrimSpace(description))
	adapterType = strings.ToLower(strings.TrimSpace(adapterType))

	switch {
	case strings.Contains(description, "wireless") || strings.Contains(description, "wifi") || strings.Contains(description, "wi-fi"):
		return entities.AdapterTypeWiFi
	case strings.Contains(description, "ethernet") || strings.Contains(adapterType, "ethernet"):
		return entities.AdapterTypeEthernet
	case strings.Contains(description, "bluetooth"):
		return entities.AdapterTypeBluetooth
	case strings.Contains(description, "virtual") || strings.Contains(description, "vmware") || strings.Contains(description, "hyper-v") || strings.Contains(description, "vbox"):
		return entities.AdapterTypeVirtual
	case strings.Contains(description, "loopback"):
		return entities.AdapterTypeLoopback
	default:
		return entities.AdapterTypeOther
	}
}

func (s *WindowsNetworkAdapterService) parseConnectionStatus(status uint32) entities.NetworkStatus {
	switch status {
	case 2:
		return entities.NetworkStatusConnected
	case 7:
		return entities.NetworkStatusDisconnected
	case 8:
		return entities.NetworkStatusConnecting
	case 9:
		return entities.NetworkStatusDisconnecting
	default:
		return entities.NetworkStatusUnknown
	}
}

func (s *WindowsNetworkAdapterService) getBasicAdapterInfo(ctx context.Context, adapterID string) (*entities.NetworkAdapter, error) {
	adapters, err := s.GetNetworkAdapters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapters: %w", err)
	}

	for _, adapter := range adapters {
		if adapter.ID == adapterID {
			return adapter, nil
		}
	}

	return nil, fmt.Errorf("adapter not found: %s", adapterID)
}

func (s *WindowsNetworkAdapterService) runElevatedCommand(ctx context.Context, name string, args ...string) error {
	// Validação de segurança para comandos
	allowedCommands := map[string]bool{
		"netsh":    true,
		"ipconfig": true,
	}
	
	if !allowedCommands[strings.ToLower(name)] {
		return fmt.Errorf("command not allowed: %s", name)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	// Define timeout para o comando
	timer := time.AfterFunc(30*time.Second, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	defer timer.Stop()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func (s *WindowsNetworkAdapterService) isThrottlingDisabled() bool {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, 
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Multimedia\SystemProfile\Tasks\Games`, 
		registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	val, _, err := key.GetIntegerValue("NetworkThrottlingIndex")
	return err == nil && val == 0xffffffff
}

func (s *WindowsNetworkAdapterService) isHighPrioritySet() bool {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, 
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Multimedia\SystemProfile\Tasks\Games`, 
		registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	val, _, err := key.GetIntegerValue("SystemResponsiveness")
	return err == nil && val == 0
}

func (s *WindowsNetworkAdapterService) isTCPOptimized() bool {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, 
		`SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, 
		registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	tcpNoDelay, _, err := key.GetIntegerValue("TCPNoDelay")
	if err != nil || tcpNoDelay != 1 {
		return false
	}

	tcpAckFreq, _, err := key.GetIntegerValue("TcpAckFrequency")
	return err == nil && tcpAckFreq == 1
}
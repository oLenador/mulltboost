package entities

import (
	"fmt"
	"time"
)

// NetworkAdapterType define os tipos de adaptadores de rede
type NetworkAdapterType int

const (
	AdapterTypeUnknown NetworkAdapterType = iota
	AdapterTypeEthernet
	AdapterTypeWiFi
	AdapterTypeBluetooth
	AdapterTypeVirtual
	AdapterTypeLoopback
	AdapterTypeOther
)

func (t NetworkAdapterType) String() string {
	switch t {
	case AdapterTypeEthernet:
		return "Ethernet"
	case AdapterTypeWiFi:
		return "WiFi"
	case AdapterTypeBluetooth:
		return "Bluetooth"
	case AdapterTypeVirtual:
		return "Virtual"
	case AdapterTypeLoopback:
		return "Loopback"
	case AdapterTypeOther:
		return "Other"
	default:
		return "Unknown"
	}
}

// NetworkStatus define os status de conexão
type NetworkStatus int

const (
	NetworkStatusUnknown NetworkStatus = iota
	NetworkStatusConnected
	NetworkStatusDisconnected
	NetworkStatusConnecting
	NetworkStatusDisconnecting
	NetworkStatusError
)

func (s NetworkStatus) String() string {
	switch s {
	case NetworkStatusConnected:
		return "Connected"
	case NetworkStatusDisconnected:
		return "Disconnected"
	case NetworkStatusConnecting:
		return "Connecting"
	case NetworkStatusDisconnecting:
		return "Disconnecting"
	case NetworkStatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// NetworkMetrics contém métricas de rede
type NetworkMetrics struct {
	BytesSent       uint64    `json:"bytes_sent"`
	BytesReceived   uint64    `json:"bytes_received"`
	PacketsSent     uint64    `json:"packets_sent"`
	PacketsReceived uint64    `json:"packets_received"`
	Errors          uint64    `json:"errors"`
	Dropped         uint64    `json:"dropped"`
	LastUpdated     time.Time `json:"last_updated"`
}

// NetworkAdapter representa um adaptador de rede do sistema
type NetworkAdapter struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	MACAddress  string             `json:"mac_address"`
	Type        NetworkAdapterType `json:"type"`
	Status      NetworkStatus      `json:"status"`
	Speed       uint64             `json:"speed"` // em bits por segundo
	IsEnabled   bool               `json:"is_enabled"`
	Metrics     *NetworkMetrics    `json:"metrics,omitempty"`
}

// IsActive verifica se o adaptador está ativo
func (na *NetworkAdapter) IsActive() bool {
	return na.IsEnabled && na.Status == NetworkStatusConnected
}

// GetSpeedMbps retorna a velocidade em Mbps
func (na *NetworkAdapter) GetSpeedMbps() float64 {
	if na.Speed == 0 {
		return 0
	}
	return float64(na.Speed) / 1_000_000 // converte de bps para Mbps
}

// Validate valida os dados do adaptador
func (na *NetworkAdapter) Validate() error {
	if na.ID == "" {
		return fmt.Errorf("adapter ID cannot be empty")
	}
	if na.Name == "" {
		return fmt.Errorf("adapter name cannot be empty")
	}
	return nil
}
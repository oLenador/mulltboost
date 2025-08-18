//go:build windows
// +build windows
package system

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"golang.org/x/sys/windows"
)

// Windows API constants and structures
const (
	DXGI_ERROR_NOT_FOUND = 0x887A0002
)

var (
	dxgi                = windows.NewLazySystemDLL("dxgi.dll")
	procCreateDXGIFactory = dxgi.NewProc("CreateDXGIFactory")
)

// DXGI structures
type DXGI_ADAPTER_DESC struct {
	Description           [128]uint16
	VendorId              uint32
	DeviceId              uint32
	SubSysId              uint32
	Revision              uint32
	DedicatedVideoMemory  uintptr
	DedicatedSystemMemory uintptr
	SharedSystemMemory    uintptr
	AdapterLuid           int64
}

// GPUInfoRepository implementation for Windows
type gpuInfoRepository struct {
	mu           sync.RWMutex
	gpuCache     []*entities.GPUInfo
	lastRefresh  time.Time
	refreshTTL   time.Duration
}

// NewGPUInfoRepository creates a new Windows GPU info repository
func NewGPUInfoRepository() *gpuInfoRepository {
	return &gpuInfoRepository{
		refreshTTL: 30 * time.Second, // Cache por 30 segundos
	}
}

// GetGPUInfo retorna informações de todas as GPUs
func (r *gpuInfoRepository) GetGPUInfo(ctx context.Context) ([]*entities.GPUInfo, error) {
	r.mu.RLock()
	if time.Since(r.lastRefresh) < r.refreshTTL && len(r.gpuCache) > 0 {
		result := make([]*entities.GPUInfo, len(r.gpuCache))
		copy(result, r.gpuCache)
		r.mu.RUnlock()
		return result, nil
	}
	r.mu.RUnlock()

	return r.refreshAndGetGPUInfo(ctx)
}

// GetPrimaryGPU retorna a GPU primária
func (r *gpuInfoRepository) GetPrimaryGPU(ctx context.Context) (*entities.GPUInfo, error) {
	gpus, err := r.GetGPUInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get GPU info: %w", err)
	}

	for _, gpu := range gpus {
		if gpu.IsPrimary {
			return gpu, nil
		}
	}

	// Se nenhuma foi marcada como primária, retorna a primeira
	if len(gpus) > 0 {
		return gpus[0], nil
	}

	return nil, fmt.Errorf("no GPU found")
}

// RefreshGPUInfo força uma atualização das informações de GPU
func (r *gpuInfoRepository) RefreshGPUInfo(ctx context.Context) error {
	_, err := r.refreshAndGetGPUInfo(ctx)
	return err
}

// GetGPUUsage retorna o uso atual da GPU
func (r *gpuInfoRepository) GetGPUUsage(ctx context.Context, gpuID string) (float64, error) {
	return r.getGPUMetricFromWMI(ctx, gpuID, "LoadPercentage")
}

// GetGPUTemperature retorna a temperatura da GPU
func (r *gpuInfoRepository) GetGPUTemperature(ctx context.Context, gpuID string) (float64, error) {
	return r.getGPUMetricFromWMI(ctx, gpuID, "Temperature")
}

// GetVRAMUsage retorna o uso de VRAM da GPU
func (r *gpuInfoRepository) GetVRAMUsage(ctx context.Context, gpuID string) (float64, error) {
	gpu, err := r.findGPUByID(ctx, gpuID)
	if err != nil {
		return 0, err
	}

	if gpu.VRAMSize == 0 {
		return 0, nil
	}

	return (float64(gpu.VRAMUsed) / float64(gpu.VRAMSize)) * 100, nil
}

// refreshAndGetGPUInfo atualiza o cache de informações de GPU
func (r *gpuInfoRepository) refreshAndGetGPUInfo(ctx context.Context) ([]*entities.GPUInfo, error) {
	// Inicializar COM
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize COM: %w", err)
	}
	defer ole.CoUninitialize()

	// Obter informações via DXGI
	dxgiGPUs, err := r.getDXGIGPUs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get DXGI GPU info: %w", err)
	}

	// Enriquecer com informações do WMI
	wmiGPUs, err := r.getWMIGPUs(ctx)
	if err != nil {
		// Log do erro mas não falha completamente
		// Em produção, usar um logger adequado
		fmt.Printf("Warning: failed to get WMI GPU info: %v\n", err)
	}

	// Combinar informações
	gpus := r.mergeGPUInfo(dxgiGPUs, wmiGPUs)

	// Determinar GPU primária
	r.setPrimaryGPU(gpus)

	// Atualizar cache
	r.mu.Lock()
	r.gpuCache = gpus
	r.lastRefresh = time.Now()
	r.mu.Unlock()

	result := make([]*entities.GPUInfo, len(gpus))
	copy(result, gpus)

	return result, nil
}

// getDXGIGPUs obtém informações de GPU via DXGI
func (r *gpuInfoRepository) getDXGIGPUs(ctx context.Context) ([]*entities.GPUInfo, error) {
	var gpus []*entities.GPUInfo

	// Criar factory DXGI
	var factory uintptr
	ret, _, _ := procCreateDXGIFactory.Call(
		uintptr(unsafe.Pointer(&IID_IDXGIFactory{})),
		uintptr(unsafe.Pointer(&factory)),
	)

	if ret != 0 {
		return nil, fmt.Errorf("failed to create DXGI factory: %x", ret)
	}

	defer windows.Syscall(factory, 0, 0, 0, 0) // Release

	// Enumerar adaptadores
	for i := uint32(0); ; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var adapter uintptr
		ret, _, _ := windows.Syscall(
			factory+3*8, // EnumAdapters offset
			3,
			uintptr(i),
			uintptr(unsafe.Pointer(&adapter)),
			0,
		)

		if ret == DXGI_ERROR_NOT_FOUND {
			break
		}
		if ret != 0 {
			continue
		}

		gpu, err := r.createGPUFromDXGIAdapter(adapter, i)
		if err != nil {
			windows.Syscall(adapter, 0, 0, 0, 0) // Release
			continue
		}

		gpus = append(gpus, gpu)
		windows.Syscall(adapter, 0, 0, 0, 0) // Release
	}

	return gpus, nil
}

// createGPUFromDXGIAdapter cria uma entidade GPUInfo a partir de um adaptador DXGI
func (r *gpuInfoRepository) createGPUFromDXGIAdapter(adapter uintptr, index uint32) (*entities.GPUInfo, error) {
	var desc DXGI_ADAPTER_DESC

	// GetDesc
	ret, _, _ := windows.Syscall(
		adapter+4*8, // GetDesc offset
		2,
		uintptr(unsafe.Pointer(&desc)),
		0,
	)

	if ret != 0 {
		return nil, fmt.Errorf("failed to get adapter description: %x", ret)
	}

	// Converter descrição UTF-16 para string
	name := windows.UTF16ToString(desc.Description[:])
	
	// Determinar vendor
	vendor := r.getVendorName(desc.VendorId)

	now := time.Now()

	gpu := &entities.GPUInfo{
		ID:               fmt.Sprintf("GPU_%d", index),
		Name:             strings.TrimSpace(name),
		Vendor:           vendor,
		DeviceID:         fmt.Sprintf("0x%04X", desc.DeviceId),
		VRAMSize:         int64(desc.DedicatedVideoMemory),
		IsDiscrete:       desc.DedicatedVideoMemory > 0,
		SupportsDirectX:  "12.0", // Assumir suporte básico
		SupportsVulkan:   true,   // Assumir suporte básico
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return gpu, nil
}

// getWMIGPUs obtém informações complementares via WMI
func (r *gpuInfoRepository) getWMIGPUs(ctx context.Context) (map[string]*entities.GPUInfo, error) {
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, fmt.Errorf("failed to create WMI locator: %w", err)
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to get WMI interface: %w", err)
	}
	defer wmi.Release()

	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer", nil, "root\\cimv2")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WMI service: %w", err)
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// Query Win32_VideoController
	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", 
		"SELECT * FROM Win32_VideoController WHERE VideoProcessor IS NOT NULL")
	if err != nil {
		return nil, fmt.Errorf("failed to execute WMI query: %w", err)
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, err := oleutil.GetProperty(result, "Count")
	if err != nil {
		return nil, fmt.Errorf("failed to get result count: %w", err)
	}

	count := int(countVar.Val)
	gpus := make(map[string]*entities.GPUInfo)

	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
		if err != nil {
			continue
		}
		item := itemRaw.ToIDispatch()

		gpu := r.createGPUFromWMIItem(item, i)
		if gpu != nil {
			gpus[gpu.Name] = gpu
		}

		item.Release()
	}

	return gpus, nil
}

// createGPUFromWMIItem cria informações de GPU a partir de um item WMI
func (r *gpuInfoRepository) createGPUFromWMIItem(item *ole.IDispatch, index int) *entities.GPUInfo {
	name := r.getWMIStringProperty(item, "Name")
	if name == "" {
		return nil
	}

	now := time.Now()

	gpu := &entities.GPUInfo{
		ID:            fmt.Sprintf("WMI_%d", index),
		Name:          name,
		DriverVersion: r.getWMIStringProperty(item, "DriverVersion"),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// AdapterRAM pode fornecer informações de VRAM
	if ramStr := r.getWMIStringProperty(item, "AdapterRAM"); ramStr != "" {
		if ram, err := strconv.ParseInt(ramStr, 10, 64); err == nil {
			gpu.VRAMSize = ram
		}
	}

	return gpu
}

// mergeGPUInfo combina informações de DXGI e WMI
func (r *gpuInfoRepository) mergeGPUInfo(dxgiGPUs []*entities.GPUInfo, wmiGPUs map[string]*entities.GPUInfo) []*entities.GPUInfo {
	for _, dxgiGPU := range dxgiGPUs {
		// Tentar encontrar correspondência no WMI pelo nome
		for _, wmiGPU := range wmiGPUs {
			if r.namesMatch(dxgiGPU.Name, wmiGPU.Name) {
				// Enriquecer informações DXGI com WMI
				if dxgiGPU.DriverVersion == "" {
					dxgiGPU.DriverVersion = wmiGPU.DriverVersion
				}
				break
			}
		}
	}

	return dxgiGPUs
}

// setPrimaryGPU determina qual GPU é primária
func (r *gpuInfoRepository) setPrimaryGPU(gpus []*entities.GPUInfo) {
	if len(gpus) == 0 {
		return
	}

	// Lógica simples: a primeira GPU discreta ou a primeira se nenhuma for discreta
	var primaryIndex = 0
	for i, gpu := range gpus {
		if gpu.IsDiscrete {
			primaryIndex = i
			break
		}
	}

	gpus[primaryIndex].IsPrimary = true
}

// Helper functions

func (r *gpuInfoRepository) getVendorName(vendorID uint32) string {
	vendors := map[uint32]string{
		0x10DE: "NVIDIA",
		0x1002: "AMD",
		0x8086: "Intel",
		0x1414: "Microsoft",
	}

	if vendor, exists := vendors[vendorID]; exists {
		return vendor
	}
	return "Unknown"
}

func (r *gpuInfoRepository) getWMIStringProperty(item *ole.IDispatch, property string) string {
	if prop, err := oleutil.GetProperty(item, property); err == nil && prop.VT == ole.VT_BSTR {
		return prop.ToString()
	}
	return ""
}

func (r *gpuInfoRepository) namesMatch(name1, name2 string) bool {
	// Lógica simples de matching de nomes
	name1 = strings.ToLower(strings.TrimSpace(name1))
	name2 = strings.ToLower(strings.TrimSpace(name2))
	
	return strings.Contains(name1, name2) || strings.Contains(name2, name1)
}

func (r *gpuInfoRepository) findGPUByID(ctx context.Context, gpuID string) (*entities.GPUInfo, error) {
	gpus, err := r.GetGPUInfo(ctx)
	if err != nil {
		return nil, err
	}

	for _, gpu := range gpus {
		if gpu.ID == gpuID {
			return gpu, nil
		}
	}

	return nil, fmt.Errorf("GPU with ID %s not found", gpuID)
}

func (r *gpuInfoRepository) getGPUMetricFromWMI(ctx context.Context, gpuID string, metric string) (float64, error) {
	// Implementação simplificada - em produção, usar queries específicas
	// para obter métricas em tempo real via WMI ou performance counters
	
	// Por enquanto, retorna valores mock baseados no contexto
	switch metric {
	case "LoadPercentage":
		return 45.5, nil // Mock value
	case "Temperature":
		return 72.0, nil // Mock value
	default:
		return 0, fmt.Errorf("unsupported metric: %s", metric)
	}
}

// IID for DXGI Factory
type IID_IDXGIFactory struct{}

var IIDDXGIFactory = ole.GUID{0x7b7166ec, 0x21c7, 0x44ae, [8]byte{0xb2, 0x1a, 0xc9, 0xae, 0x32, 0x1a, 0xe3, 0x69}}

func (IID_IDXGIFactory) GUID() *ole.GUID {
	return &IIDDXGIFactory
}
package monitoring

import (
	"context"
	"sync"
	"time"

	"github.com/oLenador/mulltbost/internal/core/application/ports/outbound"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type Service struct {
    metricsRepo outbound.SystemMetricsRepository
    monitoring  bool
    mutex       sync.RWMutex
}

func NewService(metricsRepo outbound.SystemMetricsRepository) *Service {
    return &Service{
        metricsRepo: metricsRepo,
        monitoring:  false,
    }
}

func (s *Service) GetSystemMetrics(ctx context.Context) (*entities.SystemMetrics, error) {
    cpu, err := s.metricsRepo.GetCPUMetrics(ctx)
    if err != nil {
        return nil, err
    }

    memory, err := s.metricsRepo.GetMemoryMetrics(ctx)
    if err != nil {
        return nil, err
    }

    disk, err := s.metricsRepo.GetDiskMetrics(ctx)
    if err != nil {
        return nil, err
    }

    return &entities.SystemMetrics{
        CPU:         *cpu,
        Memory:      *memory,
        Disk:        *disk,
        Timestamp:   time.Now(),
    }, nil
}

func (s *Service) StartRealTimeMonitoring(ctx context.Context, interval int) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    s.monitoring = true
    
    go func() {
        ticker := time.NewTicker(time.Duration(interval) * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                if !s.IsMonitoring() {
                    return
                }
                // Aqui você pode implementar a lógica de coleta em tempo real
                // Por exemplo, enviar métricas via WebSocket ou canal
            case <-ctx.Done():
                return
            }
        }
    }()
    
    return nil
}

func (s *Service) StopRealTimeMonitoring(ctx context.Context) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    s.monitoring = false
    return nil
}

func (s *Service) IsMonitoring() bool {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    return s.monitoring
}

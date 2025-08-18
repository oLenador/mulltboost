//go:build !windows
// +build !windows

package windows

import (
	"context"
	"log"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

// WinServiceManagerService stub para sistemas não-Windows
type WinServiceManagerService struct {
	logger *log.Logger
}

// NewWinServiceManagerService cria uma instância dummy
func NewWinServiceManagerService(logger *log.Logger) *WinServiceManagerService {
	return &WinServiceManagerService{
		logger: logger,
	}
}

func (w *WinServiceManagerService) StartService(ctx context.Context, serviceName string) error {
	w.logger.Printf("[STUB] StartService chamado em sistema não-Windows: %s", serviceName)
	return nil
}

func (w *WinServiceManagerService) StopService(ctx context.Context, serviceName string) error {
	w.logger.Printf("[STUB] StopService chamado em sistema não-Windows: %s", serviceName)
	return nil
}

func (w *WinServiceManagerService) RestartService(ctx context.Context, serviceName string) error {
	w.logger.Printf("[STUB] RestartService chamado em sistema não-Windows: %s", serviceName)
	return nil
}

func (w *WinServiceManagerService) SetServiceStartupType(ctx context.Context, serviceName, startupType string) error {
	w.logger.Printf("[STUB] SetServiceStartupType chamado em sistema não-Windows: %s -> %s", serviceName, startupType)
	return nil
}

func (w *WinServiceManagerService) GetServiceStatus(ctx context.Context, serviceName string) (*entities.WindowsService, error) {
	w.logger.Printf("[STUB] GetServiceStatus chamado em sistema não-Windows: %s", serviceName)
	return &entities.WindowsService{
		ServiceName: serviceName,
		Status:      "unknown",
	}, nil
}

func (w *WinServiceManagerService) ListServices(ctx context.Context) ([]*entities.WindowsService, error) {
	w.logger.Println("[STUB] ListServices chamado em sistema não-Windows")
	return []*entities.WindowsService{}, nil
}

func (w *WinServiceManagerService) IsServiceRunning(ctx context.Context, serviceName string) (bool, error) {
	w.logger.Printf("[STUB] IsServiceRunning chamado em sistema não-Windows: %s", serviceName)
	return false, nil
}

func (w *WinServiceManagerService) DisableUnnecessaryServices(ctx context.Context) error {
	w.logger.Println("[STUB] DisableUnnecessaryServices chamado em sistema não-Windows")
	return nil
}

func (w *WinServiceManagerService) EnableEssentialServices(ctx context.Context) error {
	w.logger.Println("[STUB] EnableEssentialServices chamado em sistema não-Windows")
	return nil
}

func (w *WinServiceManagerService) GetServiceStatistics(ctx context.Context) (map[string]interface{}, error) {
	w.logger.Println("[STUB] GetServiceStatistics chamado em sistema não-Windows")
	return map[string]interface{}{}, nil
}

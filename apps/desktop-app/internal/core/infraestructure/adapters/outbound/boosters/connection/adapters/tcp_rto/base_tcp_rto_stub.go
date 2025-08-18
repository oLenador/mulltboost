//go:build !windows
package connection

import "github.com/oLenador/mulltbost/internal/core/application/ports/inbound"

func NewTCPRTOBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
    return nil
}

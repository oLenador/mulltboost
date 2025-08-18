//go:build !windows
package connection

import "github.com/oLenador/mulltbost/internal/core/application/ports/inbound"

func NewTCPFastOpenBooster(services *inbound.ExecutorDepServices) inbound.BoosterUseCase {
    return nil
}

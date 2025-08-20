//go:build !windows && !linux
package connection

import "github.com/oLenador/mulltbost/internal/core/application/ports/inbound"

func NewDNSCacheExecutor(services *inbound.ExecutorDepServices) inbound.PlatformExecutor {
    return nil
}
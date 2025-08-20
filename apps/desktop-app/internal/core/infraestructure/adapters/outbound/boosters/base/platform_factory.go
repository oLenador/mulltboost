package booster

import "github.com/oLenador/mulltbost/internal/core/application/ports/inbound"

type ServicesFactory interface {
    CreatePlatformServices() inbound.PlatformServices
}

var DefaultFactory ServicesFactory

func GetPlatformServices() inbound.PlatformServices {
    if DefaultFactory == nil {
        panic("platform factory not initialized")
    }
    return DefaultFactory.CreatePlatformServices()
}

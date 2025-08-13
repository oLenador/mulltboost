package connection

import (
	"github.com/oLenador/mulltbost/internal/boosters/connection/boosters"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func GetAllPlugins() []inbound.BoosterUseCase {
	tcpBooster := connection.NewTCPBooster()
	dnsBooster := connection.NewDNSBooster()

	return []inbound.BoosterUseCase{
		tcpBooster,
		dnsBooster,
	}
}
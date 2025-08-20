package connection

import (
	// arpBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/arp_cache"
	// eeeBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/disable_eee"
	dnsBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/dns_cache"
	// ecnBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/ecn"
	// interruptBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/interrupt_coalescing"
	// jumboBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/jumbo_frames"
	// lsoBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/lso"
	// mtuBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/mtu_adjustment"
	// nagleBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/nagle_algorithm"
	// qosBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/qos"
	// rssBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/rss"
	// rxQueueBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/rx_queue"
	// sendBufferBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/send_buffer"
	 tcpAdvancedBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/tcp_advanced"
	 tcpCongestionBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/tcp_congestion"
	 tcpFastOpenBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/tcp_fast_open"
	 tcpRtoBooster "github.com/oLenador/mulltbost/internal/core/infraestructure/adapters/outbound/boosters/connection/adapters/tcp_rto"
	"github.com/oLenador/mulltbost/internal/core/application/ports/inbound"
)

func GetAllPlugins(services *inbound.ExecutorDepServices) []inbound.BoosterUseCase {

    boosters := []inbound.BoosterUseCase{
		dnsBooster.NewDNSBooster(services),
		// arpBooster.NewARPCacheBooster(services),
		// eeeBooster.NewEEEBooster(services),
		// ecnBooster.NewECNBooster(services),
		// interruptBooster.NewInterruptCoalescingBooster(services),
		// jumboBooster.NewJumboFramesBooster(services),
		// lsoBooster.NewJumboFramesBooster(services),
		// mtuBooster.NewMTUAdjustmentBooster(services),
		// nagleBooster.NewNagleAlgorithmBooster(services),
		// qosBooster.NewQoSBooster(services),
		// rssBooster.NewRSSBooster(services),
		// rxQueueBooster.NewRxQueueBooster(services),
		// sendBufferBooster.NewSendBufferBooster(services),
		 tcpAdvancedBooster.NewAdvancedTCPBooster(services),
		 tcpCongestionBooster.NewTCPCongestionBooster(services),
		 tcpFastOpenBooster.NewTCPFastOpenBooster(services),
		 tcpRtoBooster.NewTCPRTOBooster(services),
	}

	result := []inbound.BoosterUseCase{}
    for _, b := range boosters {
        if b != nil {
            result = append(result, b)
        }
    }

    return result
}

package connection

import (
	dnsBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/dns_cache"
	arpBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/arp_cache"
	ecnBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/ecn"
	eeeBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/disable_eee"
	interruptBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/interrupt_coalescing"
	jumboBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/jumbo_frames"
	lsoBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/lso"
	mtuBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/mtu_adjustment"
	nagleBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/nagle_algorithm"
	qosBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/qos"
	rssBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/rss"
	rxQueueBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/rx_queue"
	sendBufferBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/send_buffer"
	tcpAdvancedBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/tcp_advanced"
	tcpCongestionBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/tcp_congestion"
	tcpFastOpenBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/tcp_fast_open"
	tcpRtoBooster "github.com/oLenador/mulltbost/internal/boosters/connection/boosters/tcp_rto"

	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

func GetAllPlugins() []inbound.BoosterUseCase {
	return []inbound.BoosterUseCase{
		dnsBooster.NewDNSCacheBooster(),
		arpBooster.NewARPCacheBooster(),
		eeeBooster.NewEEEBooster(),
		ecnBooster.NewECNBooster(),
		interruptBooster.NewInterruptCoalescingBooster(),
		jumboBooster.NewJumboFramesBooster(),
		lsoBooster.NewJumboFramesBooster(),
		mtuBooster.NewMTUAdjustmentBooster(),
		nagleBooster.NewNagleAlgorithmBooster(),
		qosBooster.NewQoSBooster(),
		rssBooster.NewRSSBooster(),
		rxQueueBooster.NewRxQueueBooster(),
		sendBufferBooster.NewSendBufferBooster(),
		tcpAdvancedBooster.NewAdvancedTCPBooster(),
		tcpCongestionBooster.NewTCPCongestionBooster(),
		tcpFastOpenBooster.NewTCPFastOpenBooster(),
		tcpRtoBooster.NewTCPRTOBooster(),
	}
}

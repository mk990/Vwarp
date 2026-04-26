package warp

import (
	"math/rand"
	"net/netip"
	"time"

	"github.com/mk990/Vwarp/iputils"
)

func WarpPrefixes() []netip.Prefix {
	return []netip.Prefix{
		netip.MustParsePrefix("162.159.192.0/24"),
		netip.MustParsePrefix("162.159.195.0/24"),
		netip.MustParsePrefix("188.114.96.0/24"),
		netip.MustParsePrefix("188.114.97.0/24"),
		netip.MustParsePrefix("188.114.98.0/24"),
		netip.MustParsePrefix("188.114.99.0/24"),
		netip.MustParsePrefix("2606:4700:d0::/64"),
		netip.MustParsePrefix("2606:4700:d1::/64"),
	}
}

func RandomWarpPrefix(v4, v6 bool) netip.Prefix {
	if !v4 && !v6 {
		panic("Must choose a IP version for RandomWarpPrefix")
	}

	cidrs := WarpPrefixes()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		cidr := cidrs[rng.Intn(len(cidrs))]

		if v4 && cidr.Addr().Is4() {
			return cidr
		}

		if v6 && cidr.Addr().Is6() {
			return cidr
		}
	}
}

func WarpPorts() []uint16 {
	// Ports confirmed by the Cloudflare WARP API (in identity endpoint.ports).
	// 2408 is the primary WARP port; 500/4500 are IKE/NAT-T (used by AtomicNoize
	// obfuscation); 1701 is L2TP.
	return []uint16{
		2408,
		500,
		1701,
		4500,
	}
}

// GetWarpPorts is an alias for WarpPorts for compatibility
func GetWarpPorts() []uint16 {
	return WarpPorts()
}

func RandomWarpPort() uint16 {
	ports := WarpPorts()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ports[rng.Intn(len(ports))]
}

func RandomWarpEndpoint(v4, v6 bool) (netip.AddrPort, error) {
	randomIP, err := iputils.RandomIPFromPrefix(RandomWarpPrefix(v4, v6))
	if err != nil {
		return netip.AddrPort{}, err
	}

	return netip.AddrPortFrom(randomIP, RandomWarpPort()), nil
}

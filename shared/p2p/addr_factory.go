package p2p

import (
	"github.com/libp2p/go-libp2p/config"
	ma "github.com/multiformats/go-multiaddr"
)

// relayAddrsOnly returns an AddrFactory which will only return Multiaddr via
// specified relay string.
func relayAddrsOnly(relay string) config.AddrsFactory {
	return func(addrs []ma.Multiaddr) []ma.Multiaddr {
		if relay == "" {
			return addrs
		}

		var relayAddrs []ma.Multiaddr

		for _, a := range addrs {
			if a.String() == "/p2p-circuit" {
				continue
			}
			relayAddr, err := ma.NewMultiaddr(relay + "/p2p-circuit" + a.String())
			if err != nil {
				panic(err) // This might happen if `relay` is malformed multiaddr.
			}

			relayAddrs = append(relayAddrs, relayAddr)
		}

		return relayAddrs
	}
}

package netpubsub

import (
	"net"

	watermillnet "github.com/andyollylarkin/watermill-net"
	"github.com/andyollylarkin/watermill-net/pkg/connection"
)

func ListenerFactory(network, listenAddr string) (watermillnet.Listener, error) {
	switch network {
	case "tcp4":
		l, err := net.Listen(network, listenAddr)
		if err != nil {
			return nil, err
		}

		return connection.NewTCP4Listener(l), nil
	default:
		panic("not implemented listener")
	}
}

package netpubsub

import (
	"crypto/tls"

	watermillnet "github.com/andyollylarkin/watermill-net"
	"github.com/andyollylarkin/watermill-net/pkg/connection"
)

func ListenerFactory(network, listenAddr string) (watermillnet.Listener, error) {
	switch network {
	case "tcp4":
		return connection.NewTCP4Listener(listenAddr)
	default:
		panic("not implemented listener")
	}
}

func ListenerTlsFactory(network, listenAddr string, conf *tls.Config) (watermillnet.Listener, error) {
	switch network {
	case "tcp4":
		return connection.NewTCP4TlsListener(listenAddr, conf)
	default:
		panic("not implemented listener")
	}
}

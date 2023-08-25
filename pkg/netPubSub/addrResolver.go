package netpubsub

import "net"

func ResolveAddr(network string, addr string) (net.Addr, error) {
	switch network {
	case Tcp4:
		return net.ResolveTCPAddr(network, addr)
	default:
		panic("resolver not implemented")
	}
}

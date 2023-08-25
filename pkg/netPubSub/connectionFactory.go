package netpubsub

import (
	"context"
	"net"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	watermillnet "github.com/andyollylarkin/watermill-net"
	connection "github.com/andyollylarkin/watermill-net/pkg/connection"
)

type ReconnectionConfig struct {
	Ctx        context.Context
	Backoff    connection.Backoff
	Log        watermill.LoggerAdapter
	RemoteAddr net.Addr
	ErrFilter  connection.ErrorFilter
	RWTimeout  time.Duration
}



func ConnectionFactory(netType string, keepAlive time.Duration, rConfig *ReconnectionConfig) watermillnet.Connection {
	switch netType {
	case Tcp4:
		var conn watermillnet.Connection
		conn = connection.NewTCPConnection(net.Dialer{}, keepAlive)

		if rConfig != nil {
			conn = reconnectionWrapper(rConfig, conn)
		}

		return conn
	default:
		panic("not implemented yet")
	}
}

func reconnectionWrapper(rConfig *ReconnectionConfig, baseConn watermillnet.Connection) watermillnet.Connection {
	wConn := connection.NewReconnectWrapper(rConfig.Ctx, baseConn, rConfig.Backoff, rConfig.Log, rConfig.RemoteAddr,
		rConfig.ErrFilter, rConfig.RWTimeout)

	return wConn
}

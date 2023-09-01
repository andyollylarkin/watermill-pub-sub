package netpubsub

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	watermillnet "github.com/andyollylarkin/watermill-net"
	connection "github.com/andyollylarkin/watermill-net/pkg/connection"
)

type Config struct {
	watermillnet.PublisherConfig
	watermillnet.SubscriberConfig
	ReconnectionConfig *ReconnectionConfig
	WaitAck            bool
	KeepAlive          time.Duration
	TlsConfig          *tls.Config
	Log                watermill.LoggerAdapter
	RWTimeout          time.Duration
}

type NetPubSub struct {
	pub     *watermillnet.Publisher
	sub     *watermillnet.Subscriber
	config  Config
	started bool
}

func NewNetPubSub(config Config) (*NetPubSub, error) {
	ps := new(NetPubSub)
	ps.config = config

	p, err := watermillnet.NewPublisher(ps.config.PublisherConfig, ps.config.WaitAck)
	if err != nil {
		return nil, err
	}

	ps.pub = p

	s, err := watermillnet.NewSubscriber(ps.config.SubscriberConfig)
	if err != nil {
		return nil, err
	}

	ps.sub = s

	return ps, nil
}

// Publisher get publisher.
// If useAck is true, publisher should not wait for ack/nack message from consumer.
func (nps *NetPubSub) Publisher() message.Publisher {
	return nps.pub
}

// Subscriber get subscriber.
func (nps *NetPubSub) Subscriber() message.Subscriber {
	return nps.sub
}

// Establish bidirectional connection. Method should prepare publisher and subscriber for listen/receive.
// Blocks until the client connects.
func (nps *NetPubSub) RunAsServer(network, listerAddr string) error {
	if nps.started {
		return errors.New("net pub/sub already started")
	}

	var l watermillnet.Listener

	var err error

	if nps.config.TlsConfig != nil {
		l, err = ListenerTlsFactory(network, listerAddr, nps.config.TlsConfig)
		if err != nil {
			return err
		}
	} else {
		l, err = ListenerFactory(network, listerAddr)
		if err != nil {
			return err
		}
	}

	err = nps.sub.Connect(l)
	if err != nil {
		return err
	}

	conn, err := nps.sub.GetConnection()
	if err != nil {
		return err
	}

	conn = connection.NewReconnectListenerWrapper(context.Background(), conn, nps.config.Log, nps.config.RWTimeout, l)

	nps.pub.SetConnection(conn)

	nps.sub.SetConnection(conn)

	nps.started = true

	return nil
}

func (nps *NetPubSub) RunAsClient(network, remoteAddr string) error {
	if nps.started {
		return errors.New("net pub/sub already started")
	}

	var conn watermillnet.Connection

	if nps.config.TlsConfig != nil {
		conn = ConnectionTlsFactory(network, nps.config.KeepAlive, nps.config.ReconnectionConfig, nps.config.TlsConfig)
	} else {
		conn = ConnectionFactory(network, nps.config.KeepAlive, nps.config.ReconnectionConfig)
	}

	addr, err := ResolveAddr(network, remoteAddr)
	if err != nil {
		return err
	}

	nps.pub.SetConnection(conn)
	nps.sub.SetConnection(conn)

	err = nps.pub.Connect(addr)
	if err != nil {
		return err
	}

	err = nps.sub.Connect(nil)
	if err != nil {
		return err
	}

	nps.started = true

	return nil
}

// Close both pub and sub for communication channels.
func (nps *NetPubSub) Close() error {
	err := nps.pub.Close()
	if err != nil {
		return err
	}

	err = nps.sub.Close()
	if err != nil {
		return err
	}

	return nil
}

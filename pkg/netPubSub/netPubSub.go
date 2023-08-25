package netpubsub

import (
	"errors"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	watermillnet "github.com/andyollylarkin/watermill-net"
	watermillpubsub "github.com/andyollylarkin/watermill-pub-sub"
)

type PublisherConfig struct {
	watermillnet.PublisherConfig
	ReconnectionConfig *ReconnectionConfig
	KeepAlive          time.Duration
	WaitAck            bool
}

type SubscriberConfig struct {
	watermillnet.SubscriberConfig
}

type NetPubSub struct {
	pub              *watermillnet.Publisher
	sub              *watermillnet.Subscriber
	mode             watermillpubsub.Mode
	publisherConfig  PublisherConfig
	subscriberConfig SubscriberConfig
	started          bool
}

func NewNetPubSub(pConfig PublisherConfig, sConfig SubscriberConfig) (*NetPubSub, error) {
	ps := new(NetPubSub)
	ps.publisherConfig = pConfig
	ps.subscriberConfig = sConfig

	p, err := watermillnet.NewPublisher(ps.publisherConfig.PublisherConfig, ps.publisherConfig.WaitAck)
	if err != nil {
		return nil, err
	}

	ps.pub = p

	s, err := watermillnet.NewSubscriber(ps.subscriberConfig.SubscriberConfig)
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

	l, err := ListenerFactory(network, listerAddr)
	if err != nil {
		return err
	}

	err = nps.sub.Connect(l)
	if err != nil {
		return err
	}

	conn, err := nps.sub.GetConnection()
	if err != nil {
		return err
	}

	nps.pub.SetConnection(conn)

	nps.sub.SetConnection(conn)

	nps.started = true

	return nil
}

func (nps *NetPubSub) RunAsClient(network, remoteAddr string) error {
	if nps.started {
		return errors.New("net pub/sub already started")
	}

	conn := ConnectionFactory(network, nps.publisherConfig.KeepAlive, nps.publisherConfig.ReconnectionConfig)

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

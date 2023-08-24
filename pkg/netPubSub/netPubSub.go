package netpubsub

import (
	"github.com/ThreeDotsLabs/watermill/message"
	watermillpubsub "github.com/andyollylarkin/watermill-pub-sub"
)

type NetPubSub struct{}

// Publisher get publisher.
// If useAck is true, publisher should not wait for ack/nack message from consumer.
func (nps *NetPubSub) Publisher() message.Publisher {
	panic("not implemented") // TODO: Implement
}

// Subscriber get subscriber.
func (nps *NetPubSub) Subscriber() message.Subscriber {
	panic("not implemented") // TODO: Implement
}

// Establish bidirectional connection. Method should prepare publisher and subscriber for listen/receive.
func (nps *NetPubSub) Connect(mode watermillpubsub.Mode) error {
	panic("not implemented") // TODO: Implement
}

// Close both pub and sub for communication channels.
func (nps *NetPubSub) Close() error {
	panic("not implemented") // TODO: Implement
}

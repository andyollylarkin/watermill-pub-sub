package watermillpubsub

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

type Mode byte

const (
	ServerMode Mode = iota
	ClientMode
)

// PubSub represent bidirectional pub/sub interface for watermill publishers/subscribers.
type PubSub interface {
	// Publisher get publisher.
	Publisher() message.Publisher
	// Subscriber get subscriber.
	Subscriber() message.Subscriber
	// Establish bidirectional connection. Method should prepare publisher and subscriber for listen/receive.
	Connect(mode Mode) error
	// Close both pub and sub for communication channels.
	Close() error
}

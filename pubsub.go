package watermillpubsub

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

// PubSub represent bidirectional pub/sub interface for watermill publishers/subscribers.
type PubSub interface {
	// Publisher get publisher.
	Publisher() message.Publisher
	// Subscriber get subscriber.
	Subscriber() message.Subscriber
	// Establish bidirectional connection. Method should prepare publisher and subscriber for listen/receive.
	RunAsServer(network string, listerAddr string) error
	// Establish bidirectional connection. Method should prepare publisher and subscriber for listen/receive.
	RunAsClient(network, remoteAddr string) error
	// Close both pub and sub for communication channels.
	Close() error
}

package watermillpubsub

import (
	"context"
	"errors"
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

// Connection is bidirectional pub/sub connection.
// Create new topics for reading/writing messages in topics that are different in meaning and purpose.
type Connection struct {
	pubsubMechanism PubSub
	topics          []*Topic
}

func NewConnection(ps PubSub) *Connection {
	c := new(Connection)
	c.pubsubMechanism = ps

	return c
}

func (c *Connection) RunAsServer(network, listerAddr string) error {
	return c.pubsubMechanism.RunAsServer(network, listerAddr)
}

func (c *Connection) RunAsClient(network, remoteAddr string) error {
	return c.pubsubMechanism.RunAsClient(network, remoteAddr)
}

// Close pub/sub connection.
func (c *Connection) Close() error {
	for _, t := range c.topics {
		t.mu.Lock()
		t.closed = true
		close(t.readCh)
		t.mu.Unlock()
	}

	return c.pubsubMechanism.Close()
}

// Topic is bidirectional pub/sub topic.
type Topic struct {
	topic   string
	autoAck bool
	subCh   <-chan *message.Message
	readCh  chan *message.Message
	pub     message.Publisher
	sub     message.Subscriber
	ctx     context.Context
	mu      sync.Mutex
	closed  bool
}

func (c *Connection) NewTopic(ctx context.Context, topic string, autoAck bool) (*Topic, error) {
	t := new(Topic)
	t.autoAck = autoAck
	t.topic = topic
	t.ctx = ctx
	t.pub = c.pubsubMechanism.Publisher()
	t.sub = c.pubsubMechanism.Subscriber()

	subCh, err := t.sub.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}

	t.subCh = subCh
	t.readCh = make(chan *message.Message)

	go t.consume(subCh)

	c.topics = append(c.topics, t)

	return t, nil
}

func (t *Topic) consume(ch <-chan *message.Message) {
	for {
		select {
		case <-t.ctx.Done():
			return
		case v := <-ch:
			t.mu.Lock()
			if !t.closed {
				if t.autoAck {
					v.Ack()
				}
				t.readCh <- v
			} else {
				t.mu.Unlock()

				return
			}
			t.mu.Unlock()
		}
	}
}

// ConsumeSingle consume new single message from topic channel.
func (t *Topic) ConsumeSingle() *message.Message {
	return <-t.readCh
}

// ConsumeMulti consume new messages from topic channel.
func (t *Topic) ConsumeMulti() <-chan *message.Message {
	return t.readCh
}

// Produce produce new message to topic.
func (t *Topic) Publish(msg *message.Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return errors.New("topic closed")
	}

	if len(msg.UUID) == 0 {
		msg.UUID = watermill.NewUUID()
	}

	return t.pub.Publish(t.topic, msg)
}

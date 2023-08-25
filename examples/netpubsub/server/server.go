package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	watermillnet "github.com/andyollylarkin/watermill-net"
	wmnetPkg "github.com/andyollylarkin/watermill-net/pkg"
	watermillpubsub "github.com/andyollylarkin/watermill-pub-sub"
	netpubsub "github.com/andyollylarkin/watermill-pub-sub/pkg/netPubSub"
)

func main() {
	// logger := watermill.NewStdLogger(true, true)
	sConfig := netpubsub.SubscriberConfig{
		SubscriberConfig: watermillnet.SubscriberConfig{
			Marshaler:   wmnetPkg.MessagePackMarshaler{},
			Unmarshaler: wmnetPkg.MessagePackUnmarshaler{},
			// Logger:      logger,
		},
	}
	pConfig := netpubsub.PublisherConfig{
		KeepAlive: time.Minute * 1,
		WaitAck:   false,
		PublisherConfig: watermillnet.PublisherConfig{
			Marshaler:   wmnetPkg.MessagePackMarshaler{},
			Unmarshaler: wmnetPkg.MessagePackUnmarshaler{},
			// Logger:      logger,
		},
	}
	nps, err := netpubsub.NewNetPubSub(pConfig, sConfig)

	if err != nil {
		log.Fatal(err)
	}

	err = nps.RunAsServer(netpubsub.Tcp4, "127.0.0.1:9090")
	if err != nil {
		log.Fatal(err)
	}

	conn := watermillpubsub.NewConnection(nps)
	t, err := conn.NewTopic(context.Background(), "topic1", true)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for m := range t.ConsumeMulti() {
			fmt.Println(string(m.Payload))
		}

		fmt.Println("DONE CONSUME")
	}()

	for {
		time.Sleep(time.Second * 2)

		err = t.Publish(message.NewMessage("", message.Payload("Hello from server")))

		if err != nil {
			log.Fatal(err)
		}
	}
}

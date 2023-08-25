package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	watermillnet "github.com/andyollylarkin/watermill-net"
	wmnetPkg "github.com/andyollylarkin/watermill-net/pkg"
	"github.com/andyollylarkin/watermill-net/pkg/connection"
	watermillpubsub "github.com/andyollylarkin/watermill-pub-sub"
	netpubsub "github.com/andyollylarkin/watermill-pub-sub/pkg/netPubSub"
	"github.com/sethvargo/go-retry"
)

func main() {
	remoteAddr := &net.TCPAddr{IP: []byte{127, 0, 0, 1}, Port: 9090}
	logger := watermill.NewStdLogger(true, true)
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
		ReconnectionConfig: &netpubsub.ReconnectionConfig{
			Ctx:        context.Background(),
			Backoff:    retry.NewConstant(time.Second * 3),
			Log:        logger,
			RemoteAddr: remoteAddr,
			ErrFilter:  connection.DefaultErrorFilter,
			RWTimeout:  time.Minute * 1,
		},
	}
	nps, err := netpubsub.NewNetPubSub(pConfig, sConfig)

	if err != nil {
		log.Fatal(err)
	}

	err = nps.RunAsClient(netpubsub.Tcp4, remoteAddr.String())
	if err != nil {
		log.Fatal(err)
	}

	conn := watermillpubsub.NewConnection(nps)
	t, err := conn.NewTopic(context.Background(), "topic1", true)

	if err != nil {
		log.Fatal(err)
	}

	// server -> client
	go func() {
		for m := range t.ConsumeMulti() {
			fmt.Println(string(m.Payload))
		}

		fmt.Println("DONE CONSUME")
	}()

	// client -> server
	for {
		time.Sleep(time.Second * 2)

		err = t.Publish(message.NewMessage("", message.Payload("Hello from client")))

		if err != nil {
			log.Fatal(err)
		}
	}
}

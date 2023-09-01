package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	watermillnet "github.com/andyollylarkin/watermill-net"
	wmnetPkg "github.com/andyollylarkin/watermill-net/pkg"
	"github.com/andyollylarkin/watermill-net/pkg/connection"
	connectionhelpers "github.com/andyollylarkin/watermill-net/pkg/helpers/connectionHelpers"
	watermillpubsub "github.com/andyollylarkin/watermill-pub-sub"
	netpubsub "github.com/andyollylarkin/watermill-pub-sub/pkg/netPubSub"
	"github.com/sethvargo/go-retry"
)

func main() {
	logger := watermill.NewStdLogger(true, true)

	cert, err := connectionhelpers.LoadCerts("/home/denis/Desktop/local.test.ru.crt",
		"/home/denis/Desktop/local.test.ru.key")
	if err != nil {
		log.Fatal(err)
	}

	config := netpubsub.Config{
		PublisherConfig: watermillnet.PublisherConfig{
			Marshaler:   wmnetPkg.MessagePackMarshaler{},
			Unmarshaler: wmnetPkg.MessagePackUnmarshaler{},
			Logger:      logger,
		},
		SubscriberConfig: watermillnet.SubscriberConfig{
			Marshaler:   wmnetPkg.MessagePackMarshaler{},
			Unmarshaler: wmnetPkg.MessagePackUnmarshaler{},
			Logger:      logger,
		},
		TlsConfig: &tls.Config{
			Certificates: cert,
		},
		ReconnectionConfig: &netpubsub.ReconnectionConfig{
			Ctx:       context.Background(),
			Backoff:   retry.NewConstant(time.Second * 5),
			Log:       logger,
			ErrFilter: connection.DefaultErrorFilter,
			RWTimeout: time.Second * 5,
		},
		Log:       logger,
		RWTimeout: time.Second * 5,
	}

	nps, err := netpubsub.NewNetPubSub(config)

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

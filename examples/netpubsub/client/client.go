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

func main() { //nolint
	remoteAddr, err := netpubsub.ResolveAddr(netpubsub.Tcp4, ":9090")
	if err != nil {
		log.Fatal(err)
	}

	logger := watermill.NewStdLogger(true, true)

	cert, err := connectionhelpers.LoadCerts("/home/denis/Desktop/local.test.ru.crt",
		"/home/denis/Desktop/local.test.ru.key")
	if err != nil {
		log.Fatal(err)
	}

	rootCertPool, err := connectionhelpers.LoadCertPool("/home/denis/Documents/rootca.crt")
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
		WaitAck: false,
		ReconnectionConfig: &netpubsub.ReconnectionConfig{
			Ctx:        context.Background(),
			Backoff:    retry.NewConstant(time.Second * 3),
			Log:        logger,
			RemoteAddr: remoteAddr,
			ErrFilter:  connection.DefaultErrorFilter,
			RWTimeout:  time.Minute * 1,
		},
		TlsConfig: &tls.Config{
			Certificates: cert,
			RootCAs:      rootCertPool,
			ServerName:   "local.test.ru",
		},
	}

	nps, err := netpubsub.NewNetPubSub(config)
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

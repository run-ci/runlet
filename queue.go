package main

import (
	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
)

func SubscribeToQueue(subject string) (<-chan *nats.Msg, func()) {
	// TODO: make this configurable
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("error listening to subject %v: %v", subject, err)
	}

	ch := make(chan *nats.Msg)
	sub, err := nc.ChanSubscribe(subject, ch)
	if err != nil {
		log.Fatalf("error listening to subject %v: %v", subject, err)
	}

	teardown := func() {
		defer nc.Close()

		err := sub.Unsubscribe()
		if err != nil {
			log.Fatalf("error unsubscribing from %v: %v", subject, err)
		}

		close(ch)
	}

	return ch, teardown
}

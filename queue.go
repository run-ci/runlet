package main

import (
	"math"
	"time"

	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
)

// SubscribeToQueue sets up a subscription in NATS on the "pipelines" subject.
// TODO: abstract away the dependency on NATS.
func SubscribeToQueue(url, subject string) (<-chan *nats.Msg, func()) {
	log.Infof("connecting to nats at %v with subject %v", url, subject)

	nc, err := nats.Connect(url)
	if err != nil {
		for i := 1; i <= 3; i++ {
			timeout := time.Duration(math.Pow(2, float64(i))) * time.Second

			log.Infof("error connecting to nats: %v", err)
			log.Infof("retry %v with timeout %v", i, timeout)

			time.Sleep(timeout)
			nc, err = nats.Connect(url)
			if err == nil {
				break
			}
		}
	}

	log.Info("nats connection successful")

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

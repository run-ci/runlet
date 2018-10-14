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
	logger.Info("connecting to nats")

	nc, err := nats.Connect(url)
	if err != nil {
		for i := 1; i <= 3; i++ {
			timeout := time.Duration(math.Pow(2, float64(i))) * time.Second

			logger.WithFields(log.Fields{
				"error": err,
			}).Warnf("error connecting to nats, retrying after %v seconds", timeout)

			time.Sleep(timeout)
			nc, err = nats.Connect(url)
			if err == nil {
				break
			}
		}
	}

	logger.Info("nats connection successful")

	// TODO: ensure that this is a queue-group and not a publisher. Multiple runlets
	// shouldn't process the same message.
	ch := make(chan *nats.Msg)
	sub, err := nc.ChanSubscribe(subject, ch)
	if err != nil {
		logger.Fatalf("error listening to subject %v: %v", subject, err)
	}

	logger.Debugf("subscribed to subject %v", subject)

	teardown := func() {
		logger.Debugf("begin tearing down nats connection")

		defer nc.Close()

		err := sub.Unsubscribe()
		if err != nil {
			logger.WithFields(log.Fields{
				"error": err,
			}).Fatalf("unable to cleanly unsubscribe from subject %v", subject)
		}

		close(ch)
	}

	return ch, teardown
}

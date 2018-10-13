package main

import (
	"encoding/json"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	"gitlab.com/run-ci/run/pkg/run"
)

var natsURL string

func init() {
	// TODO: fix all this logging shit...
	log.SetLevel(log.DebugLevel)

	natsURL = os.Getenv("RUNLET_NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
}

func main() {
	log.Info("booting runlet")

	evq, teardown := SubscribeToQueue(natsURL, "pipelines")
	defer teardown()

	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatalf("error opening docker socket: %v", err)
	}

	log.Debug("initialized docker client")

	agent, err := run.NewAgent(client)
	if err != nil {
		log.Fatalf("error initializing run agent with our client: %v", err)
	}

	for msg := range evq {
		log.Debugf("processing message: %s", msg.Data)

		var ev Event
		err := json.Unmarshal(msg.Data, &ev)
		if err != nil {
			log.Warnf("error parsing message: %v", err)
		}

		log.Debugf("parsed event: %+v", ev)

		vol := initCIVolume(agent, client, ev.Remote)

		for name, tasks := range ev.Steps {
			log.Debugf("running step %v", name)

			for _, task := range tasks {
				if task.Mount == "" {
					task.Mount = "/ci/repo"
				}

				if task.Shell == "" {
					task.Shell = "sh"
				}

				spec := run.ContainerSpec{
					Imgref: task.Image,
					Cmd:    task.GetCmd(),
					Mount: run.Mount{
						Src:   vol,
						Point: task.Mount,
						Type:  "volume",
					},
				}

				id, status, err := agent.RunContainer(spec)
				if err != nil {
					log.Fatalf("error running task container: %v", err)
				}

				log.Debugf("task container %v exited with status %v", id, status)
			}
		}
	}
}

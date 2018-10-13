package main

import (
	"os"

	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
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

	// client, err := docker.NewClient("unix:///var/run/docker.sock")
	// if err != nil {
	// 	log.Fatalf("error opening docker socket: %v", err)
	// }

	// log.Debug("initialized docker client")

	// agent, err := run.NewAgent(client)
	// if err != nil {
	// 	log.Fatalf("error initializing run agent with our client: %v", err)
	// }

	for ev := range evq {
		log.Debugf("processing message: %s", ev.Data)
	}

	// vol := initCIVolume(agent, client)

	// task := run.Task{
	// 	Image:   "golang:1.11-stretch",
	// 	Mount:   "/go/src/github.com/juicemia/go-sample-app",
	// 	Shell:   "sh",
	// 	Command: "go build -v",
	// }

	// spec := run.ContainerSpec{
	// 	Imgref: task.Image,
	// 	Cmd:    task.GetCmd(),
	// 	Mount: run.Mount{
	// 		Src:   vol,
	// 		Point: task.Mount,
	// 		Type:  "volume",
	// 	},
	// }

	// id, status, err := agent.RunContainer(spec)
	// if err != nil {
	// 	log.Fatalf("error running task container: %v", err)
	// }

	// log.Debugf("task container %v exited with status %v", id, status)
}

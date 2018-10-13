package main

import (
	"gitlab.com/run-ci/run/pkg/run"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

func init() {
	// TODO: fix all this logging shit...
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("booting runlet")

	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatalf("error opening docker socket: %v", err)
	}

	log.Debug("initialized docker client")

	agent, err := run.NewAgent(client)
	if err != nil {
		log.Fatalf("error initializing run agent with our client: %v", err)
	}

	vol := initCIVolume(agent, client)

	task := run.Task{
		Image:   "golang:1.11-stretch",
		Mount:   "/go/src/github.com/juicemia/go-sample-app",
		Shell:   "sh",
		Command: "go build -v",
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

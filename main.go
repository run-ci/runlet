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

	err = agent.VerifyImagePresent("run-ci/git-clone", true)
	if err != nil {
		log.Fatalf("error verifying image git-clone image presence: %v", err)
	}

	vol, err := client.CreateVolume(docker.CreateVolumeOptions{
		Name: "runlet-test-vol",
	})
	if err != nil {
		log.Fatalf("error creating test volume: %v", err)
	}

	log.Debugf("created volume: %v", vol.Name)

	spec := run.ContainerSpec{
		Imgref: "run-ci/git-clone",
		Cmd:    []string{"https://github.com/juicemia/go-sample", "repo"},
		Mount: run.Mount{
			Src:   vol.Name,
			Point: "/ci",
			Type:  "volume",
		},
	}

	id, status, err := agent.RunContainer(spec)
	if err != nil {
		log.Fatalf("error running git clone container %v: %v", id, err)
	}

	log.Debugf("git clone container exited with status %v", status)
}

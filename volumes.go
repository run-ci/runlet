package main

import (
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
	"gitlab.com/run-ci/run/pkg/run"
)

func initCIVolume(agent *run.Agent, client *docker.Client, remote string) string {
	logger := logger.WithField("remote", remote)

	name := "runlet-test-vol"

	// TODO: don't hard-code this maybe?
	err := agent.VerifyImagePresent("run-ci/git-clone", true)
	if err != nil {
		logger.WithField("error", err).
			Fatalf("unable to verify git-clone image presence")
	}

	vol, err := client.CreateVolume(docker.CreateVolumeOptions{
		// TODO: dynamically generate this, the names shouldn't matter
		Name: name,
	})
	if err != nil {
		logger.WithField("error", err).
			Fatalf("unable to create test volume")
	}

	logger = logger.WithField("vol", vol.Name)
	logger.Debugf("created volume: %v", vol.Name)

	spec := run.ContainerSpec{
		// TODO: fix all this hard-coded crap
		Imgref: "run-ci/git-clone",
		Cmd:    []string{remote, "."},
		Mount: run.Mount{
			Src:   vol.Name,
			Point: "/ci",
			Type:  "volume",
		},
	}

	logger.Debug("populating volumes")

	id, status, err := agent.RunContainer(spec)
	if err != nil {
		log.Fatalf("error running git clone container %v: %v", id, err)
	}

	logger.Debugf("git clone container exited with status %v", status)

	return name
}

package main

import (
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
	"gitlab.com/run-ci/run/pkg/run"
)

func initCIVolume(agent *run.Agent, client *docker.Client) string {
	name := "runlet-test-vol"
	gitrepo := "https://github.com/juicemia/go-sample-app"

	// TODO: don't hard-code this maybe?
	err := agent.VerifyImagePresent("run-ci/git-clone", true)
	if err != nil {
		log.Fatalf("error verifying image git-clone image presence: %v", err)
	}

	vol, err := client.CreateVolume(docker.CreateVolumeOptions{
		// TODO: dynamically generate this, the names shouldn't matter
		Name: name,
	})
	if err != nil {
		log.Fatalf("error creating test volume: %v", err)
	}

	log.Debugf("created volume: %v", vol.Name)

	spec := run.ContainerSpec{
		// TODO: fix all this hard-coded crap
		Imgref: "run-ci/git-clone",
		Cmd:    []string{gitrepo, "."},
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

	return name
}

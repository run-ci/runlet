package main

import (
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("booting runlet")

	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatalf("error opening docker socket: %v", err)
	}

	log.Debug("initialized docker client")

	err = initCIVolume(client)
	if err != nil {
		log.Fatalf("error initializing CI volume: %v", err)
	}
}

func initCIVolume(client *docker.Client) error {
	vol, err := client.CreateVolume(docker.CreateVolumeOptions{
		Name: "runlet-test-vol",
	})
	if err != nil {
		return err
	}

	log.Debugf("created volume: %v", vol.Name)

	imgs, err := client.ListImages(docker.ListImagesOptions{
		All:    true,
		Filter: "run-ci/git-clone",
	})
	if err != nil {
		log.Fatalf("error searching for image %v: %v", "run-ci/git-clone", err)
	}

	if len(imgs) < 0 {
		log.Debugf("image %v not found, pulling", "run-ci/git-clone")

		pullopts := docker.PullImageOptions{
			Repository: "run-ci/git-clone",
		}

		pullopts.OutputStream = os.Stdout

		err = client.PullImage(pullopts, docker.AuthConfiguration{})
		if err != nil {
			log.Fatalf("error pulling image %v: %v", "run-ci/git-clone", err)
		}

		log.Debugf("image %v pulled", "run-ci/git-clone")
	}

	ccfg := &docker.Config{
		Image:        "run-ci/git-clone",
		Cmd:          []string{"https://github.com/juicemia/go-sample", "repo"},
		AttachStderr: true,
		AttachStdout: true,
		Volumes: map[string]struct{}{
			"/mnt/git-clone": struct{}{},
		},
		WorkingDir: "/mnt/git-clone",
	}
	hcfg := &docker.HostConfig{
		Mounts: []docker.HostMount{
			docker.HostMount{
				Target: "/mnt/git-clone",
				Source: "runlet-test-vol",
				Type:   "volume",
			},

			docker.HostMount{
				Target: "/var/run/docker.sock",
				Source: "/var/run/docker.sock",
				Type:   "bind",
			},
		},
	}

	ncfg := &docker.NetworkingConfig{}

	cnt, err := client.CreateContainer(docker.CreateContainerOptions{
		Config:           ccfg,
		HostConfig:       hcfg,
		NetworkingConfig: ncfg,
	})
	if err != nil {
		log.Fatalf("error creating container for task: %v", err)
	}

	log.Debugf("container %v created", cnt.ID)

	go func() {
		log.Debugf("attaching container %v", cnt.ID)

		attachcfg := docker.AttachToContainerOptions{
			Container: cnt.ID,
			Stderr:    true,
			Stdout:    true,
			Stream:    true,
			Logs:      true,

			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
		}

		err = client.AttachToContainer(attachcfg)
		if err != nil {
			msgs := fmt.Sprintf("error attaching to task container: %v", err)

			err := cleanupContainer(client, cnt.ID)
			if err != nil {
				msgs = fmt.Sprintf("%v\nerror cleaning up container %v: %v", msgs, cnt.ID, err)
			}

			log.Fatalf(msgs)
		}
	}()

	log.Debugf("starting container %v", cnt.ID)

	err = client.StartContainer(cnt.ID, cnt.HostConfig)
	if err != nil {
		msgs := fmt.Sprintf("error starting task container: %v", err)

		err := cleanupContainer(client, cnt.ID)
		if err != nil {
			msgs = fmt.Sprintf("%v\nerror cleaning up container %v: %v", msgs, cnt.ID, err)
		}

		log.Fatalf(msgs)
	}

	status, err := client.WaitContainer(cnt.ID)
	if err != nil {
		msgs := fmt.Sprintf("error running task container: %v", err)

		err := cleanupContainer(client, cnt.ID)
		if err != nil {
			msgs = fmt.Sprintf("%v\nerror cleaning up container %v: %v", msgs, cnt.ID, err)
		}

		log.Fatalf(msgs)
	}

	fmt.Printf("task container exited with status %v\n", status)

	err = cleanupContainer(client, cnt.ID)
	if err != nil {
		log.Fatalf("error cleaning up container %v: %v", cnt.ID, err)
	}
	return nil
}

func cleanupCIVolume(client *docker.Client, name string) error {
	return client.RemoveVolume(name)
}

func cleanupContainer(client *docker.Client, id string) error {
	return client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            id,
		RemoveVolumes: true,
		Force:         true,
	})
}

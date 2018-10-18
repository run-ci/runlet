package main

import (
	"encoding/json"
	"os"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	nats "github.com/nats-io/go-nats"
	"github.com/run-ci/run/pkg/run"
	"github.com/run-ci/runlet/store"
	log "github.com/sirupsen/logrus"
)

var natsURL, gitimg, cimnt string
var logger *log.Entry

func init() {
	natsURL = os.Getenv("RUNLET_NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	gitimg = os.Getenv("RUNLET_GIT_IMAGE")
	if gitimg == "" {
		gitimg = "run-ci/git-clone"
	}

	cimnt = os.Getenv("RUNLET_CI_MOUNT")
	if cimnt == "" {
		cimnt = "/ci/repo"
	}

	switch strings.ToLower(os.Getenv("RUNLET_LOG_LEVEL")) {
	case "debug", "trace":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	logger = log.WithFields(log.Fields{
		"package": "main",
	})
}

func main() {
	logger.Info("booting runlet")

	evq, teardown := SubscribeToQueue(natsURL, "pipelines", "runlet")
	defer teardown()

	st, err := store.NewPostgres("postgres://runlet_test:runlet_test@store:5432/runlet_test?sslmode=disable")
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to postgres")
	}

	logger.Info("connecting to database")

	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		logger.WithFields(log.Fields{
			"error": err,
		}).Fatal("error opening docker socket")
	}

	logger.Info("initialized docker client")

	agent, err := run.NewAgent(client)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to initialize run agent")
	}

	logger.Info("initialized run agent")

	for msg := range evq {
		logger.Debugf("processing message %s", msg.Data)

		var ev Event
		err := json.Unmarshal(msg.Data, &ev)
		if err != nil {
			logger.WithFields(log.Fields{
				"error": err,
			}).Error("error parsing message, skipping")

			continue
		}

		vol := initCIVolume(agent, client, ev.Remote)

		p := &store.Pipeline{
			Remote: ev.Remote,
		}

		err = st.LoadPipeline(p)
		if err != nil {
			logger.WithField("error", err).Error("error loading pipeline from store, skipping")

			continue
		}

		for name, tasks := range ev.Steps {
			logger := logger.WithFields(log.Fields{
				"remote": ev.Remote,
				"step":   name,
			})

			logger.Debug("running step")

			for _, task := range tasks {
				logger := logger.WithField("task", task.Name)

				if task.Mount == "" {
					task.Mount = cimnt

					logger.Debugf("mount point set to %v", task.Mount)
				}

				if task.Shell == "" {
					task.Shell = "sh"

					logger.Debugf("shell set to %v", task.Shell)
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
				logger = logger.WithField("container_id", id)
				if err != nil {
					logger.WithField("error", err).
						Fatalf("error running task container")
				}

				logger.Debugf("task container exited with status %v", status)
			}
		}

		err = client.RemoveVolume(vol)
		if err != nil {
			logger.WithFields(log.Fields{
				"error": err,
				"vol":   vol,
			}).Error("unable to delete volume")
		}
	}
}

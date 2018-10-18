package store

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var logger *log.Entry

func init() {
	logger = log.WithFields(log.Fields{
		"package": "store",
	})
}

type PipelineStore interface {
	SavePipeline(Pipeline) error
	LoadPipeline(*Pipeline) error
}

type Pipeline struct {
	Remote string `db:"remote"`
	Name   string `db:"name"`
	Ref    string `db:"ref"`
	Runs   []Run  `db:"-"`
}

type Run struct {
	Count   int        `db:"id"`
	Start   *time.Time `db:"start"`
	End     *time.Time `db:"end"`
	Success bool       `db:"success"`
	Steps   []Step     `db:"-"`

	PipelineRemote string `db:"pipeline_remote"`
	PipelineName   string `db:"pipeline_name"`
}

type Step struct {
	ID      int        `db:"id"`
	Name    string     `db:"name"`
	Start   *time.Time `db:"start"`
	End     *time.Time `db:"end"`
	Tasks   []Task     `db:"-"`
	Success bool       `db:"success"`

	PipelineRemote string `db:"pipeline_remote"`
	PipelineName   string `db:"pipeline_name"`
	RunCount       int    `db:"run_count"`
}

type Task struct {
	ID      int        `db:"id"`
	Name    string     `db:"name"`
	Start   *time.Time `db:"start"`
	End     *time.Time `db:"end"`
	Success bool       `db:"success"`
}

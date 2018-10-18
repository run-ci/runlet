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
	Runs   []int  `db:"runs"`
}

type Run struct {
	ID      int        `db:"id"`
	Start   *time.Time `db:"start"`
	End     *time.Time `db:"end"`
	Success bool       `db:"success"`
	Steps   []int      `db:"steps"`
}

type Step struct {
	ID      int    `db:"id"`
	Name    string `db:"name"`
	Tasks   []int  `db:"tasks"`
	Success bool   `db:"success"`
}

type Task struct {
	ID      int        `db:"id"`
	Name    string     `db:"name"`
	Start   *time.Time `db:"start"`
	End     *time.Time `db:"end"`
	Success bool       `db:"success"`
}

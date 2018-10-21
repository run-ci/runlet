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
	ReadPipeline(*Pipeline) error

	CreateRun(*Run) error
	CreateStep(*Step) error
	CreateTask(*Task) error

	UpdateRun(*Run) error
	UpdateStep(*Step) error
	UpdateTask(*Task) error
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
	Success *bool      `db:"success"` // mid-run is neither success nor failure
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
	Success *bool      `db:"success"` // mid-run is neither success nor failure

	PipelineRemote string `db:"pipeline_remote"`
	PipelineName   string `db:"pipeline_name"`
	RunCount       int    `db:"run_count"`
}

type Task struct {
	ID      int        `db:"id"`
	Name    string     `db:"name"`
	Start   *time.Time `db:"start"`
	End     *time.Time `db:"end"`
	Success *bool      `db:"success"` // mid-run is neither success nor failure

	StepID int `db:"step_id"`
}

func (r *Run) SetStart() {
	t := time.Now()
	r.Start = &t
}

func (r *Run) SetEnd() {
	t := time.Now()
	r.End = &t
}

func (r *Run) MarkSuccess(s bool) {
	r.Success = &s
}

func (r *Run) Failed() bool {
	return r.Success != nil && *r.Success == false
}

func (st *Step) SetStart() {
	t := time.Now()
	st.Start = &t
}

func (st *Step) SetEnd() {
	t := time.Now()
	st.End = &t
}

func (st *Step) MarkSuccess(s bool) {
	st.Success = &s
}

func (task *Task) SetStart() {
	t := time.Now()
	task.Start = &t
}

func (task *Task) SetEnd() {
	t := time.Now()
	task.End = &t
}

func (task *Task) MarkSuccess(s bool) {
	task.Success = &s
}

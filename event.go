package main

import (
	"gitlab.com/run-ci/run/pkg/run"
)

// Event is a message that comes in requesting a pipeline run.
type Event struct {
	Remote string `json:"remote"`
	Steps  Steps  `json:"steps"`
}

// Steps is a list of mappings between a name and a group of
// tasks to run.
//
// TODO: make this match what's in gitlab.com/run-ci/webhooks/pkg
type Steps map[string][]run.Task

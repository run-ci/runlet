package main

import "gitlab.com/run-ci/run/pkg/run"

// Event is a message that comes in requesting a pipeline run.
type Event struct {
	Remote string              `json:"remote"`
	Nodes  map[string]run.Task `json:"nodes"`
	Edges  map[string]bool     `json:"edges"`
}

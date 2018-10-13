package main

import (
	"errors"

	"gitlab.com/run-ci/run/pkg/run"
)

// Event is a message that comes in requesting a pipeline run.
type Event struct {
	Remote string              `json:"remote"`
	Nodes  map[string]run.Task `json:"nodes"`
	Edges  map[string][]string `json:"edges"`
}

// Cycle-check
//
// DFS, marking every node that's visited
// Pipeline graphs can have more than one "root", so
// create a "fake" root that's a parent of all the
// real roots.

// CycleCheck uses a depth-first search to check for cycles in
// the event's pipeline graph. If the check fails a non-nil error
// is returned.
func (ev Event) CycleCheck() error {
	// find the "roots"

	// depth first search

	return errors.New("not implemented")
}

func (ev Event) getRoots() []string {
	// get a list of all targets in the edges map
	// roots are nodes that aren't in that list

	ret := []string{}
	tmp := map[string]struct{}{}

	for _, nodeEdges := range ev.Edges {
		for _, e := range nodeEdges {
			tmp[e] = struct{}{}
		}
	}

	for n := range ev.Nodes {
		if _, ok := tmp[n]; !ok {
			ret = append(ret, n)
		}
	}

	return ret
}

package main

import (
	"testing"

	"gitlab.com/run-ci/run/pkg/run"
)

func TestEventCycleCheck(t *testing.T) {
	ev := Event{
		Nodes: map[string]run.Task{
			"abc": run.Task{
				Name: "foo",
			},
			"def": run.Task{
				Name: "bar",
			},
		},
	}

	err := ev.CycleCheck()
	if err != nil {
		t.Logf("event: %+v", ev)
		t.Error("expected cycle check to pass")
	}
}

func TestGetRoots(t *testing.T) {
	cases := []struct {
		name     string
		ev       Event
		expected []string
	}{
		{
			"simple",

			Event{
				Nodes: map[string]run.Task{
					"abc": run.Task{
						Name: "foo",
					},
					"def": run.Task{
						Name: "bar",
					},
				},

				Edges: map[string][]string{
					"abc": []string{"def"},
				},
			},

			[]string{
				"abc",
			},
		},

		{
			"two_roots",

			Event{
				Nodes: map[string]run.Task{
					"abc": run.Task{
						Name: "foo",
					},
					"def": run.Task{
						Name: "bar",
					},
					"012": run.Task{
						Name: "baz",
					},
				},

				Edges: map[string][]string{
					"abc": []string{"012"},
					"def": []string{"012"},
				},
			},

			[]string{
				"abc",
				"def",
			},
		},
	}

	for _, cas := range cases {
		t.Run(cas.name, func(t *testing.T) {
			actual := cas.ev.getRoots()
			if len(actual) != len(cas.expected) {
				t.Logf("event: %+v", cas.ev)
				t.Errorf("expected roots: %v\n\ngot: %v\n", cas.expected, actual)
			}

			tmp := map[string]struct{}{}
			for _, n := range actual {
				tmp[n] = struct{}{}
			}

			for _, n := range cas.expected {
				if _, ok := tmp[n]; !ok {
					t.Logf("event: %+v", cas.ev)
					t.Errorf("expected to find %v in roots", n)
				}
			}

			tmp = map[string]struct{}{}
			for _, n := range cas.expected {
				tmp[n] = struct{}{}
			}

			for _, n := range actual {
				if _, ok := tmp[n]; !ok {
					t.Logf("event: %+v", cas.ev)
					t.Errorf("didn't to find %v in roots", n)
				}
			}
		})
	}
}

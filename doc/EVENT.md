# Events

Events are what drive the runlet to do work.

They must include everything the runlet needs to run
a pipeline. For now, the list of minimum things is:

1. The remote URL to clone from Git.
2. All the tasks that can possibly be run.
3. The list of connections between the tasks.

Pipelines are directed, acyclical graphs. Their nodes are tasks and
their edges dictate the sequence of events in the pipeline. As a result,
the bulk of an event is the graph describing the pipeline.

The edges are an adjacency list, but implemented as a map. This is because
since the graph is directed and reverse traversal is to be avoided, there's
no need to remember where the traversal was before, only to know what's
coming up next. The list itself is implemented as a map-backed set to make
it easy to look up edges, which is useful for checking for acyclicality.

Events must describe a whole entire pipeline because pipelines must be run
on the same Docker host. If that restriction gets lifted, this schema will
probably need to be changed substantially.

For now, they are specified as JSON, though other serialization types
may eventually be supported.

## Event Schema

Generated using [this](https://easy-json-schema.github.io/) and subsequently
edited by hand.

```JSON
{
  "type": "object",
  "required": ["remote", "nodes", "edges"],
  "properties": {
    "remote": {
      "type": "string"
    },
    "nodes": {
      "type": "object",
      "required": [],
      "patternProperties": {
        ".*": {
          "type": "object",
          "required": ["task", "image", "command"],
          "properties": {
            "task": { "type": "string" },
            "image": { "type": "string" },
            "command": { "type": "string" },
            "mount": { "type": "string" },
            "arguments": { "type": "object" }
          }
        }
      }
    },
    "edges": {
      "type": "object",
      "required": [],
      "patternProperties": {
        ".*": {
          "type": "object",
          "required": [],
          "patternProperties": {
            ".*": {
              "type": "boolean"
            }
          }
        }
      }
    }
  }
}
```

## Sample Event

```JSON
{
    "remote": "https://github.com/juicemia/go-sample.git",
    "nodes": {
        "123": {
            "task": "lint",
            "image": "golang:1.11-alpine",
            "mount": "/go/src/github.com/juicemia/go-sample",
            "command": "./scripts/go-linter.sh",
            "arguments": {
                "GOOS": "linux"
            }
        },
        "456": {
            "task": "build",
            "image": "golang:1.11-alpine",
            "mount": "/go/src/github.com/juicemia/go-sample",
            "command": "go build -v",
            "arguments": {
                "GOOS": "linux"
            }
        },
        "789": {
            "task": "test",
            "image": "golang:1.11-alpine",
            "mount": "/go/src/github.com/juicemia/go-sample",
            "command": "go test -v ./...",
            "arguments": {
                "GOOS": "linux"
            }
        },
        "ABC": {
            "task": "deploy-release-candidate",
            "image": "ruby:2.4.1-stretch",
            "command": "bundle exec rake deploy",
            "arguments": {
                "VERSION": "0.0.1"
            }
        },
        "DEF": {
            "task": "integrate",
            "image": "ruby:2.4.1-stretch",
            "command": "bundle exec rake spec"
        },
        "0AF": {
            "task": "release",
            "image": "ruby:2.4.1-stretch",
            "command": "bundle exec rake release",
            "arguments": {
                "VERSION": "0.0.1"
            }
        }
    },
    "edges": {
        "123": {
            "ABC": true
        },
        "456": {
            "ABC": true
        },
        "789": {
            "ABC": true
        },
        "ABC": {
            "DEF": true
        },
        "DEF": {
            "0AF": true
        }
    }
}
```
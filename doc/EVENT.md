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
  "required": ["remote", "steps"],
  "properties": {
    "remote": {
      "type": "string"
    },
    "steps": {
      "type": "object",
      "patternProperties": {
        ".*": {
          "type": "object",
          "required": ["tasks"],
          "properties": {
            "next": {
              "type": "string"
            },
            "tasks": {
              "type": "array",
              "items": {
                "type": "object",
                "required": ["name", "image", "command"],
                "properties": {
                  "name": { "type": "string" },
                  "image": { "type": "string" },
                  "command": { "type": "string" },
                  "mount": { "type": "string" },
                  "arguments": { "type": "object" }
                }
              }
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
  "remote": "https://github.com/juicemia/go-sample-app",

  "steps": {
    "test": {
      "tasks": [
        {
          "name": "fmt",
          "command": "./scripts/checkfmt.sh",
          "image": "golang:1.11-stretch"
        },
        {
          "name": "test",
          "command": "go test -v ./...",
          "image": "golang:1.11-stretch"
        }
      ],

      "next": "build"
    },

    "build": {
      "tasks": [
        {
          "name": "build",
          "command": "go build -v",
          "image": "golang:1.11-stretch"
        }
      ]
    }
  }
}
```
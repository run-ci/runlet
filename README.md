# runlet

Agent for `run` pipelines on Docker hosts.

## Developing

**Always** source an environment from the `env` directory.
For local development `env/local` works fine.

```bash
. env/local
```

To get a working `runlet` running on your dev workstation:

```bash
run build # build a binary

docker-compose down && docker-compose up --build
```

This spins up NATS along with a `runlet` instance.

To test it, install a NATS client. The NATS ruby gem works fine.

```bash
gem install nats
```

With the `nats` gem, trigger pipeline runs like so:

```bash
# $TEST_EVENT is defined in the previously sourced `env` file. Make
# sure there are quotes around it here because it's not all on one line. 
nats-pub pipelines runlet "$TEST_EVENT"
```

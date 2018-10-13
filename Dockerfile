FROM ubuntu:18.04

ADD runlet /bin/runlet

ENTRYPOINT ["/bin/runlet"]

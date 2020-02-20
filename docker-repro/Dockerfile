# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.12-alpine

ENV GO111MODULE=on
WORKDIR /work
COPY sample.go /work

RUN ls -l
RUN CGO_ENABLED=0 GOOS=linux go build sample.go

# Run the web service on container startup.
ENTRYPOINT ["/work/sample"]
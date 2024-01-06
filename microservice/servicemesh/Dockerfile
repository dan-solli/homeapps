# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build-stage

WORKDIR /app

COPY *.go database service x509 ./
COPY ../../common ./
COPY ../../proto ./
RUN go mod init github.com/dan-solli/homeapps/microservice/servicemesh

RUN go CGO_ENABLED=0 GOOS=linux build -o /servicemesh

FROM build-stage AS run-test-stage
RUN go test -v ./... -race -coverprofile=coverage.out -covermode=atomic
RUN go test -v ./... -tag=integration -race -coverprofile=coverage.out -covermode=atomic
RUN go test -v ./... -tag=load -race -coverprofile=coverage.out -covermode=atomic

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /servicemesh /servicemesh 

ARG VERSION=0.0.1

USER root

ENTRYPOINT ["/servicemesh"]

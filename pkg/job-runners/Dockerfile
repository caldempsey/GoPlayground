# Build the server from the Golang image as base, copy to Amazon linux Instance.
FROM golang:1.12-alpine AS build
# Instruction to run from Gomod vendor
ENV GO111MODULE auto

COPY . /src/github.com/mmacheerpuppy/GoPlayground/pkg/etl-runnner
WORKDIR /src/github.com/mmacheerpuppy/GoPlayground/pkg/etl-runnner
RUN apk add --no-cache git
# Enable vendoring mode to support common libraries.
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o _build/etl-runner .

# Extract binary from source.
# Setup Amazon Linux as maximally compatible managed environment for AWS deployment.
FROM amazonlinux:2
RUN mkdir /var/task
WORKDIR /var/task
COPY --from=build /src/github.com/mmacheerpuppy/GoPlayground/pkg/etl-runnner/_build /var/task
CMD ./etl-runnner

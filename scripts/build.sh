#!/bin/bash
set -xe

VERSION=${VERSION:-`cat VERSION.TXT`}

docker_build() {
  # Build the binary for alpine linux
  docker run --rm -it \
    -v `pwd`:/go/src/github.com/jspeyside/alarmclock speyside/golang \
    'set -x && cd /go/src/github.com/jspeyside/alarmclock && \
     export VERSION=`cat VERSION.txt` && go get ./... && \
     go build -o alarmclock -ldflags "-X github.com/jspeyside/alarmclock/domain.Version=$$VERSION"'

  # Build and tag the image
  docker build -t speyside/alarmclock:`cat VERSION.txt` .
  docker tag speyside/alarmclock:`cat VERSION.txt` speyside/alarmclock:latest
}

binary_build() {
  build_os=(darwin linux windows)

  # Build the binary for each os
  for os in ${build_os[@]}; do
      GOOS=$os GOARCH=amd64 go build -o build/alarmclock_${os} -ldflags "-X github.com/jspeyside/alarmclock/domain.Version=${VERSION}"
  done

}

mkdir -p build
docker_build
binary_build

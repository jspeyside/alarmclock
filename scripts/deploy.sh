#!/bin/bash
set -e

b2 authorize-account $B2_ID $B2_KEY
VERSION=${VERSION:-`cat VERSION.TXT`}

go get -t -u -v ./...
mkdir -p build
rm -rf build/*

# linux amd64
GOOS=linux GOARCH=amd64 go build -o build/alarmclock -ldflags "-X github.com/jspeyside/alarmclock/domain.Version=${VERSION}"
tar -czvf alarmclock_${VERSION}.tar.gz alarmclock
b2 upload-file boomerain-web alarmclock_${VERSION}.tar.gz binaries/alarmclock/linux/amd64/alarmclock_${VERSION}.tar.gz
rm -rf build/*

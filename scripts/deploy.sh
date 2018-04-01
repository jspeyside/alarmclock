#!/bin/bash
set -e

b2 authorize-account $B2_ID $B2_KEY
VERSION=${VERSION:-`cat VERSION.TXT`}

if (! docker images | grep -q speyside/alarmclock | grep -q $VERSION); then
  scripts/build.sh
fi

go get -t -u -v ./...
mkdir -p build
rm -rf build/*

# Build the binaries
cd build/
build_os=(darwin linux windows)
for os in ${build_os[@]}; do
  cp alarmclock_${os} alarmclock
  tar -czvf alarmclock_${VERSION}_${os}.tar.gz alarmclock
  b2 upload-file boomerain-web alarmclock_${VERSION}.tar.gz binaries/alarmclock/${VERSION}/alarmclock_${VERSION}_${os}.tar.gz
done
rm -rf build/*

docker push speyside/alarmclock:$VERSION speyside/alarmclock:latest

#!/bin/bash
set -e

b2 authorize-account $B2_ID $B2_KEY
VERSION=${VERSION:-`cat VERSION`}

if (! docker images | grep speyside/alarmclock | grep -q $VERSION); then
  scripts/build.sh
fi

# Build the binaries
cd build/
build_os=(darwin linux windows)
for os in ${build_os[@]}; do
  cp alarmclock_${os} alarmclock
  tar -czvf alarmclock_${VERSION}_${os}.tar.gz alarmclock
  b2 upload-file boomerain-web alarmclock_${VERSION}_${os}.tar.gz binaries/alarmclock/${VERSION}/alarmclock_${VERSION}_${os}.tar.gz
done
cd ..
rm -rf build/*

docker push speyside/alarmclock:$VERSION
docker push speyside/alarmclock:latest

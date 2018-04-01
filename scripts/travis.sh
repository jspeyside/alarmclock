#!/bin/bash
set -e

pip install --user -U b2

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
scripts/build.sh
scripts/deploy.sh

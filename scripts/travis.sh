#!/bin/bash
set -e

pip install -U pip b2

scripts/build.sh
scripts/deploy.sh

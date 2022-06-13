#!/usr/bin/env bash

set -euo pipefail

current=$(pwd)

for folder in `find . -name '*.feature' | sed -r 's|/[^/]+$||' | sort -u`; do
    (cd ${folder} && godog run *.feature)
done

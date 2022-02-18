#!/usr/bin/env bash

set -euo pipefail

IMAGE_DIR=bin

(cd proto && buf mod update)
# buf breaking --against '.git#branch=main'
buf lint
mkdir -p ${IMAGE_DIR}
buf build -o ${IMAGE_DIR}/image.bin
buf generate

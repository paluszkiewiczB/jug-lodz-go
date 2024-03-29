#!/bin/bash

set -e
set -o pipefail

docker build -t 01_hello:latest .
docker run 01_hello:latest

# inspect content of the image using dive
# https://github.com/wagoodman/dive?tab=readme-ov-file#installation
#
# dive 01_hello:latest
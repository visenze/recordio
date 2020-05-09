#! /bin/bash
BUILD_IMAGE="visenze/golang:1.10"

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
REPO_ROOT=$(dirname "$SCRIPT_DIR")

CONTAINER_REPO_ROOT="/go/src/github.com/visenze/recordio"
echo "executing $SCRIPT_DIR/docker_build.sh"
echo "repo root = $REPO_ROOT"

echo "mounting $REPO_ROOT --> $CONTAINER_REPO_ROOT"

docker run --rm -it\
    -v "$REPO_ROOT":"$CONTAINER_REPO_ROOT" \
    -w "$CONTAINER_REPO_ROOT" \
    "$BUILD_IMAGE" bash  -c "dep init 2>/dev/null; cd python && ./build.sh"

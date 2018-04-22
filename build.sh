#!/usr/bin/env bash

set -e

TARGETS=(
  'linux amd64'
  'darwin amd64'
  'windows amd64 .exe'
)

mkdir -p bin
rm -f bin/lwc-*

for tgt in "${TARGETS[@]}"; do
  set -- $tgt
  out="bin/lwc-$1-$2$3"
  echo "Building $out ..."
  env GOOS=$1 GOARCH=$2 go build -ldflags '-s -w' -o "$out" ./cmd/lwc
done

echo Done.

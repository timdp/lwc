#!/bin/bash

set -e

ver="v$1"

if git rev-parse $ver >/dev/null 2>&1; then
  echo Error: tag $ver already exists >&2
  exit 1
fi

git commit --allow-empty -m "$ver"
git tag -a "$ver" -m "$ver"

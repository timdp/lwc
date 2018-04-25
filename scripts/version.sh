#!/bin/bash

set -e

ver="v$1"

git commit --allow-empty -m "$ver"
git tag -a "$ver" -m "$ver"

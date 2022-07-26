#!/bin/bash

set -e

[ -t 1 ] && TTY=1 || TTY=0
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
LWC="$ROOT/bin/lwc"
WC=$( which wc )
FLAGS='-l -w -m'
# Detect GNU
wc -L </dev/null >/dev/null 2>&1 && FLAGS="$FLAGS -c -L"

failures=0
for path in test/fixtures/*; do
  name="$( basename "$path" )"
  [[ $TTY = 1 ]] && echo -n "[????] $name"
  set -- $( "$WC" $FLAGS <"$path" )
  wc_counts=$@
  set -- $( "$LWC" $FLAGS <"$path" )
  lwc_counts=$@
  [[ $TTY = 1 ]] && echo -ne "\r"
  if [[ $wc_counts = $lwc_counts ]]; then
    echo "[PASS] $name"
  else
    echo "[FAIL] $name"
    echo "       wc:  $wc_counts"
    echo "       lwc: $lwc_counts"
    (( failures++ ))
  fi
done

exit $failures

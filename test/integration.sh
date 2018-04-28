#!/bin/bash

set -e

[ -t 1 ] && TTY=1 || TTY=0
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
LWC="$ROOT/bin/lwc-debug"
WC="$( which wc )"

failures=0
for path in test/fixtures/*; do
  name="$( basename "$path" )"
  [[ $TTY = 1 ]] && echo -n "[????] $name"
  set -- $( "$WC" -l -w -m -c <"$path" )
  wc_counts=$@
  set -- $( "$LWC" -l -w -m -c <"$path" )
  lwc_counts=$@
  [[ $TTY = 1 ]] && echo -ne "\r"
  if [[ $wc_counts = $lwc_counts ]]; then
    echo "[PASS] $name"
  else
    echo "[FAIL] $name"
    echo "       Expected: $wc_counts"
    echo "       Actual:   $lwc_counts"
    (( failures++ ))
  fi
done

exit $failures

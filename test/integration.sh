#!/bin/bash

set -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
LWC="$ROOT/bin/lwc-debug"
WC="$( which wc )"

failures=0
for path in test/fixtures/*; do
  name="$( basename "$path" )"
  echo -n "[????] $name"
  set -- $( "$WC" -l -w -m -c <"$path" )
  wc_counts=$@
  set -- $( "$LWC" -l -w -m -c <"$path" )
  lwc_counts=$@
  if [[ $wc_counts = $lwc_counts ]]; then
    echo -e "\r[PASS] $name"
  else
    echo -e "\r[FAIL] $name"
    echo "       Expected: $wc_counts"
    echo "       Actual:   $lwc_counts"
    (( failures++ ))
  fi
done

exit $failures

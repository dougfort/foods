#! /bin/bash
# run the binary for the foods server

set -e
set -x

$GOPATH/bin/foods 2> /tmp/foods.log

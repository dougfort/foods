#! /bin/bash
# run the binary for the foods server

set -e
set -x

$GOPATH/bin/foods \
	--token-path=$GOPATH/src/github.com/dougfort/foods/tokens.json \
	2> /tmp/foods.log

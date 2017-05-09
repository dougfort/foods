#! /bin/bash
# run the binary for the foods server

set -e
set -x

$GOPATH/bin/testclient --make-tokens \
	--token-path=$GOPATH/src/github.com/dougfort/foods/tokens.json

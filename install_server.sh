#! /bin/bash
# create the binary for the foods server at $GOPATH/bin/foods

set -e
set -x

pushd $GOPATH/src/github.com/dougfort/foods
go install -race
popd

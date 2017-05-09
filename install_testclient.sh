#! /bin/bash
# create the binary for the foods test client at $GOPATH/bin/testclient

set -e
set -x

pushd $GOPATH/src/github.com/dougfort/foods/testclient
go install
popd

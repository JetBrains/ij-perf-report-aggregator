#!/bin/bash

set -ex

# go bug - even if module is replaced to local dir, still, network request is performed and proxy maybe not yet aware about module
#export GOPROXY=https://proxy.golang.org

cd ./pkg
go get -u ./...
go mod tidy
cd ..

for appDir in `find ./cmd/* -type d`
do
  printf "\nUpdate deps for $appDir\n"
  cd "$appDir"
  go get -u ./...
  go mod tidy
  cd ../..
done
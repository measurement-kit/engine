#!/bin/sh
set -ex

go get -v ./internal/cmd/mkengine-version

pkg="github.com/measurement-kit/engine/mobile/android/mkengine"
v=`${GOPATH}/bin/mkengine-version`
engine=mkengine-${v}.aar

time gomobile bind      \
  -target android       \
  -o ${engine}          \
  -ldflags="-s -w"      \
  -tags="android ${1}"  \
  ${pkg}

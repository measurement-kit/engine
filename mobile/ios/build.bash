#!/bin/sh
set -ex

go get -v ./internal/cmd/mkengine-version

pkg="github.com/measurement-kit/engine/mobile/ios/MKEngine"
v=`${GOPATH}/bin/mkengine-version`
engine=MKEngine.framework

time gomobile bind  \
  -target ios       \
  -o ${engine}      \
  -ldflags="-s -w"  \
  -tags="ios ${1}"  \
  ${pkg}

tarball=MKEngine-${v}.framework.tar
tar -cvf ${tarball} ${engine}
gzip -9 ${tarball}

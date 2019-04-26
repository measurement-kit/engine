#!/bin/sh
set -ex
gomobile bind -target android -o mkengine.aar -javapkg io.ooni.mk \
  github.com/measurement-kit/engine
gomobile bind -target ios -o MKEngine.framework -prefix MK \
  github.com/measurement-kit/engine

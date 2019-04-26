#!/bin/sh
set -ex

gomobile bind -target android -o mkengine.aar -javapkg io.ooni.mk.engine       \
  github.com/measurement-kit/engine/collector                                  \
  github.com/measurement-kit/engine/geoip                                      \
  github.com/measurement-kit/engine/task

gomobile bind -target ios -o MKEngine.framework -prefix MKE                    \
  github.com/measurement-kit/engine/collector                                  \
  github.com/measurement-kit/engine/geoip                                      \
  github.com/measurement-kit/engine/task

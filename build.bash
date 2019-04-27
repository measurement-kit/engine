#!/bin/sh
set -ex

engine="github.com/measurement-kit/engine"
pkgs=`echo ${engine}/{collector,geoip,task}`

#
# Android
#

time gomobile bind                             \
  -target android                              \
  -o mkengine.aar                              \
  -javapkg io.ooni.mk.engine                   \
  -ldflags="-s -w"                             \
  ${pkgs}

#
# iOS
#

time gomobile bind                             \
  -target ios                                  \
  -o MKEngine.framework                        \
  -prefix MKE                                  \
  -ldflags="-s -w"                             \
  ${pkgs}

versionA=MKEngine.framework/Versions/A

mv ${versionA}/{Collector,MKEngine}

cat ${versionA}/Headers/Collector.h                               \
  | sed 's/__Collector_FRAMEWORK_H__/__MKEngine_FRAMEWORK_H__/g'  \
    > ${versionA}/Headers/MKEngine.h
rm ${versionA}/Headers/Collector.h

cat ${versionA}/Modules/module.modulemap \
  | sed 's/Collector/MKEngine/g' > ${versionA}/Modules/module.modulemap.new
mv ${versionA}/Modules/module.modulemap{.new,}

rm MKEngine.framework/Collector
ln -s Versions/Current/MKEngine MKEngine.framework/MKEngine

tar -cvzf MKEngine.framework.tgz MKEngine.framework

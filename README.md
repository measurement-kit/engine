[![GoDoc](https://godoc.org/github.com/measurement-kit/engine?status.svg)](https://godoc.org/github.com/measurement-kit/engine) [![Build Status](https://travis-ci.org/measurement-kit/engine.svg?branch=master)](https://travis-ci.org/measurement-kit/engine) [![Coverage Status](https://coveralls.io/repos/github/measurement-kit/engine/badge.svg?branch=master)](https://coveralls.io/github/measurement-kit/engine?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/measurement-kit/engine)](https://goreportcard.com/report/github.com/measurement-kit/engine)

# Measurement Kit engine written in Go

This repository contains a Measurement Kit engine written in Go. The design
is described in [DESIGN.md](DESIGN.md).

You can easily integrate this repository into your Go code as usual by
adding a

```Go
import "github.com/measurement-kit/engine"
```

statement to your Go code. Make sure you have set the `GOPATH` environment
variable, or use Go modules.

You can build an AAR for Android using

```bash
./mobile/android/build.bash
```

which will regenerate the Android specific code and use `gomobile`
to generate bindings and an APK for Android. Make sure you have
the Android SDK and NDK installed, that you have run `gomobile init`,
and that you have exported `ANDROID_HOME` to point to the place in
which the Android SDK is installed.

In a similar fashion

```bash
./mobile/ios/build.bash
```

will generate a `MKEngine.framework` for iOS devices. For this to
work, you must be running macOS and have Xcode installed.

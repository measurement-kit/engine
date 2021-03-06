# Measurement Kit Engine in Go

| Author       | Simone Basso
|--------------|--------------|
| Last-Updated | 2019-05-07   |
| Status       | Pending      |

## Abstract

This is a new engine for Measurement Kit (MK) written in Go. It is motivated
by the desire to integrate a [Psiphon](https://www.psiphon3.com) network
test (nettest). Because Psiphon is written in Go, integrating
it implies integrating the Go runtime. In turn, this fact opens up
the opportunity of rewriting parts of MK in Go. In fact,
adding more code to the Go engine (e.g. adding support for talking to
OONI backend servers) would not increase much the binary size. At the
same time, a Go codebase is generally easier to deal with, and also
faster to cross compile, than a C++ codebase. This document describes
how we plan on introducing this new engine into the MK and OONI ecosystems,
and further discuss possible issues, such as code bloat, and how we
could eventually fully replace the C++ implementation.

## Background

We determined that we do not want to have the Go implementation
call code in the C++ implementation. Likewise, we do not want the C++
implementation to depend on the Go implementation. Both approaches
have been discarded, because they entail some complexity and some
maintenance toil. In particular, in both cases we would need to write
(or auto-generate) code for bridging Go and C++. Moreover, having
the C++ code depend on the Go code means that we could not take advantage
of the automatic bindings for Android and iOS generated using [Go
mobile](https://godoc.org/golang.org/x/mobile/cmd/gomobile). Additionally,
this approach would require non-OONI MK integrators to include the
Go engine even if they do not want to use its functionality, with a
corresponding increase in code size.

Thus, we concluded that the optimal plan is this:

1. for [ooni/probe-ios](https://github.com/ooni/probe-ios), we will use
`gomobile bind` to generate a framework to be integrated by OONI in addition
to the already existing framework generated using
[mkall-ios](https://github.com/measurement-kit/mkall-ios);

2. for [ooni/probe-android](https://github.com/ooni/probe-android), we
will use `gomobile bind` in a similar fashion, and OONI for Android
will include both the gomobile-generated AAR and the AAR generated using
[android-libs](https://github.com/measurement-kit/android-libs);

3. for [ooni/probe-cli](https://github.com/ooni/probe-cli) it does not
matter much what we choose to do, because both probe-cli and the MK Go API
at [go-measurement-kit](
https://github.com/measurement-kit/go-measurement-kit) are clearly
labelled as beta quality software and there are no production deployments;

4. other MK integrators must not be indirectly impacted by these changes
as long as the C++ implementation is maintained (as of this writing we
don't have any plans on deprecating it).

This strategy is sound because the two implementations can cohexist
side by side in OONI apps for quite some time. At the beginning, OONI
will use the Go code for Psiphon, and the C++ code for all
the other functionality. The apps will be more bloated, but—as mentioned
above—we want to run Psiphon tests, so we cannot avoid that. This is
why we need to discuss code bloat issues as part of this document.

In going forward, we anticipate some pressure towards
rewriting more and more MK functionality in Go. Our
current experience shows that a Go codebase is more easily testable
and maintainable than one in C++.
Therefore, we need a mechanism by which we can incrementally and
transparently switch over to the Go implementation.

One last aspect to discuss is a plan for deprecating and replacing the
C++ implementation, in case we decide the Go implementation is good
and we actually want to stop using C++. What needs to happen in this case,
in particular, is a mean for C/C++ users of the C++ engine to switch
over to the Go implementation somehow.

## Limitations

The current version of this document does not address the proper
way of dealing with users of the `measurement_kit` binary. A future
version of this design document will also discuss this case.

## ooni/probe-ios changes

We will write Go code in measurement-kit/engine such that the
automatically generated ObjectiveC bindings have the same
API signatures of the code in mkall-ios. The main differences
would be (1) that code generated from measurement-kit/engine will
use the `MKEngine` prefix rather than the `MK` prefix and (2)
that the code generated from Go may contain more methods.

For example, given this Go API

```Go
type MKEngineSomeTask struct {
  Timeout int64
  URL string
}
```

The auto-generated bindings will contain both getters and
setters for `Timeout` and `URL`, while the corresponding
exiting bindings currently only contain setters for attributes
that only need to be set. We investigated and noticed that,
if we only wante to have setter, we'd need to write:

```Go
type MKEngineSomeTask struct {
  timeout int64
  url string
}

func (st *MKEngineSomeTask) SetTimeout(v int64) {
  st.timeout = v
}

func (st *MKEngineSomeTask) SetURL(v string) {
  st.url = v
}
```

However, the above is not only more verbose but also more
complex to deal with, because we need to write different
code for iOS—where getters are named like `foo`—and Android—
where getters are named like `getFoo`. Hence, not having to
dealing with manually specifying getters and setters to make
the API more similar to the existing API also brings with
it additional maintenance toil that we'd like to avoid.

The OONI app will link with both the mkall-ios and the
measurement-kit/engine frameworks. We will choose the
implementation to use by refactoring the code and changing
the prefix used to invoke specific APIs.

All the existing non-OONI users of mkall-ios won't be
affected. Also, if they want to take advantage of
the Go engine, they can do what OONI does.

## ooni/probe-android changes

Android code generated by android-libs uses the
`io.ooni.mk` package and classes use the `MK` prefix. For
code autogenerated by measurement-kit/engine, we will
use the `io.ooni.mk.engine` package. Except from that, the code generated
from measurement-kit/engine will use an API similar to
the one generated by android-libs. Also, as discussed
for iOS, the API generated from Go will contain more
getters and setters than the android-libs API.

The OONI app will link both APKs and we will choose the
implementation to use by refactoring the code. All the existing
users of the android-libs won't be affected. Also, if they
want, they can copy OONI and use the Go engine as well.

## ooni/probe-cli changes

Any functionality not currently exposed by go-measurement-kit
will be implemented directly in measurement-kit/engine as a
public API and ooni/probe-cli will use it directly.

We will add to go-measurement-kit a new package that allows to
transparently execute either nettests implemented in C++ or
nettests implemented in Go. Because this functionality is being
added as a new package, existing go-measurement-kit users will
not see any changes. The ooni/probe-cli app, instead, will be
refactored to use this new package and see the changes.

## Code bloat discussion

Psiphon is implemented in Go and ships relatively small apps,
therefore it is possible somehow to address the code bloat. It
would be smart to ask Psiphon develpers for help.

We discussed this topic during a recent OONI community meeting
and we concluded it was also smart to do user research to
understand how the app size is an issue for our users. Another
immediate outcome of the meeting was that a 50 MiB app would
most likely not be acceptable.

In addition, here are some possible ideas for addressing the
code bloat that came out from some research.

### Android

It seems reasonable to recommend apps to implement
[APK splits](
https://developer.android.com/studio/build/configure-apk-splits) such
that different ABIs are packaged as different APKs. Another aspect
to consider is that we can implement on-demand resources downloading
rather than shipping them along with the APKs.

Regarding the APK split, a build of measurement-kit/engine that only
includes report resubmission for all ABIs yields a 8.9 MiB APK. For
comparison, an arm64-only build is 2.2 MiB.

The uncompressed cumulative size of the MMDB databases is 5 MiB. If we
can download them on demand, when they are first requested by an app, it
means we can ship OONI apps that are 5 MiB smaller.

### iOS

As for Android, we can avoid shipping the MMDB databases by default, thus
making the apps 5 MiB smaller. In particular, one interesting option to
explore on iOS are [on-demand resources](https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/On_Demand_Resources_Guide/index.html),
i.e. resources that are hosted separately from the main app bundle and
downloaded from the store on demand. Another option would be to reuse
code implemented for Android and download the resources ourselves.

It is also worth further investigation [understanding whether app
thinning requires bitcode](https://forums.developer.apple.com/thread/14650). On
this note it is also important to keep in mind that Go and other languages
are making progress towards [uploading binaries that allow the App store to
perform bitcode processing even though the binary is not bitcode](
https://github.com/golang/go/issues/22395). This latter aspect does not
directly impact apps size, but potentially allows an app using only MK's
Go engine to work on TvOS and watchOS in the future.

## On replacing the C++ implementation

We currently don't plan on deprecating the C++ implementation. Yet a
robust design document for the Go engine must also address this.

The only API that is officially exposed by the C++ implementation is
[the FFI task API allowing to run async tasks](
https://github.com/measurement-kit/measurement-kit/tree/87d63ced0fb5e74db5e6568ce1871057a4d5a43a/include/measurement_kit#measurement-kit-api).
Currently only nettests can be run as async tasks. It is also
worth mentioning that such API is a pure C API consisting of the
following functions:

```C
#include <measurement_kit/ffi.h>

typedef          struct mk_event_   mk_event_t;
typedef          struct mk_task_    mk_task_t;

mk_task_t       *mk_task_start(const char *settings);
mk_event_t      *mk_task_wait_for_next_event(mk_task_t *task);
int              mk_task_is_done(mk_task_t *task);
void             mk_task_interrupt(mk_task_t *task);

const char      *mk_event_serialization(mk_event_t *event);
void             mk_event_destroy(mk_event_t *event);

void             mk_task_destroy(mk_task_t *task);
```

This API is simple enough to rewrite in Go using
[CGO](https://golang.org/cmd/cgo/). We can return
handles rather than pointers to C users, [as
described by this article](
http://justinfx.com/2016/05/14/cpp-bindings-for-go/). This
will not change much from the ABI perspective because we'll
still be returning a `uintptr_t` value to C users, except
that such value would a index into a mutex protected
table managed by Go, rather than a pointer to real memory.

Thus, it seems doable to write this CGO code and then
maybe we can write a super tiny C library that makes sure
the exposed API really looks like C. Since C users are
mainly going to be on desktop platforms, it seems we can
get away with assembling a complete library by using

```
go build -buildmode=c-archive
```

and then linking this archive with the small C shim
into a dynamic library.

## Architecture and implementation

The top-level Go package in this repository should implement
an API such that, with minimal changes, we can obtain with
`gomobile` an API as close as possible to the Android and iOS one.

The requirement that the code should work with `gomobile` implies
that factory functions should return pointers to structures. Since
our iOS and Android API uses getters and setters, it stems that
also the toplevel Go API should also follow that style, rather than
directly exposing fields to mobile users.

The implementation of the top-level Go package should be as small
as possible and ideally just defer to internal packages that
implement the actual operations. For example:

```Go
package engine

import (
	"github.com/measurement-kit/engine/internal/collector"
)

// CollectorSubmitResults contains the results of submitting or resubmitting
// a measurement to the OONI collector.
type CollectorSubmitResults struct {
	collector.SubmitResults
}

// Good returns whether we succeeded or not.
func (r *CollectorSubmitResults) Good() bool {
	return r.SubmitResults.Good
}

// Logs returns logs useful for debugging.
func (r *CollectorSubmitResults) Logs() string {
	return r.SubmitResults.Logs
}

// UpdatedReportID returns the updated report ID.
func (r *CollectorSubmitResults) UpdatedReportID() string {
	return r.SubmitResults.UpdatedReportID
}

// UpdatedSerializedMeasurement returns the measurement with updated fields.
func (r *CollectorSubmitResults) UpdatedSerializedMeasurement() string {
	return r.SubmitResults.UpdatedSerializedMeasurement
}

// CollectorSubmitTask is a synchronous task for submitting or resubmitting a
// specific OONI measurement to the OONI collector.
type CollectorSubmitTask struct {
	collector.SubmitTask
}

const defaultTimeout int64 = 30 // seconds

// NewCollectorSubmitTask creates a new CollectorSubmitTask with the specified
// software name, software version, and serialized measurement fields.
func NewCollectorSubmitTask(swName, swVersion, measurement string) *CollectorSubmitTask {
	return &CollectorSubmitTask{
		collector.SubmitTask{
			SerializedMeasurement: measurement,
			SoftwareName:          swName,
			SoftwareVersion:       swVersion,
			Timeout:               defaultTimeout,
		},
	}
}

// SetSerializedMeasurement sets the measurement to submit.
func (t *CollectorSubmitTask) SetSerializedMeasurement(measurement string) {
	t.SubmitTask.SerializedMeasurement = measurement
}

// SetSoftwareName sets the name of the software submitting the measurement.
func (t *CollectorSubmitTask) SetSoftwareName(softwareName string) {
	t.SubmitTask.SoftwareName = softwareName
}

// SetSoftwareVersion sets the name of the software submitting the measurement.
func (t *CollectorSubmitTask) SetSoftwareVersion(softwareVersion string) {
	t.SubmitTask.SoftwareVersion = softwareVersion
}

// SetTimeout sets the number of seconds after which we abort submitting.
func (t *CollectorSubmitTask) SetTimeout(timeout int64) {
	t.SubmitTask.Timeout = timeout
}

// Run submits (or resubmits) a measurement and returns the results.
func (t *CollectorSubmitTask) Run() *CollectorSubmitResults {
	return &CollectorSubmitResults{
		SubmitResults: t.SubmitTask.Run(),
	}
}
```

This is also the API that the ooni/probe-cli should use. In specific
cases, we may make additional APIs required by probe-cli public. However,
trying to keep a consistent API everywhere is probably ideal.

The procedure that generates mobile code should ideally be such that
it starts from the above Go code with mimimal changes. In particular,
the following changes are required:

1. consistently with Effective Go suggestions, the Go code should
not use `Get` in front of getters. However, in Java a getter is
expected to start with `get`. Therefore, when generating Go mobile
code, we should add `Get` in front of all getters.

2. the package name for iOS needs to be `MKEngine` such that
we generate the correct framework name and headers. This
modification is quite easy to implement.

The above changes could either be implemented with scripts that
edit inline the source code. As of merging this document into
master, this is the strategy that we are using.

## Conclusion

We will write measurement-kit/engine in Go in such a way that
the autogenerated bindings for Android and iOS will have the
same API of existing Android and iOS libs. The names of the classes
and packages, however, will be different. OONI apps, and other
interested parties, will simply refactor their code when they want
to switch from the C++ to the Go implementation.

We will add a package to go-measurement-kit that automatically
uses either the C++ or the Go implementation depending on the
requested nettest. OONI for desktop will use this new package and
existing users will use go-measurement-kit, thus being not
affected by these changes.

In practice, we will write a minimal Go shim implementing a
suitable Go API that generates the desired bindings. All the
code that is not needed to be publish must be `internal`.

If we ever choose to stop supporting the C++ engine, there is
a simple way with which we can provide to its C/C++ users a
replacement library with a C API.

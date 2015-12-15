grpcinstrument
==============

This is a golang package that helps instrumenting [gRPC][] servers. It includes
a [Protocol Buffer][pb] compiler plugin for generating code that wraps and
instruments gRPC servers and an implementation of the `Instrumentator` interface
that logs requests to any of the backends supported by [protolog][] and exposes
metrics about requests using [Prometheus][].

[Checkout the godoc][godoc] online or download this package and view them
locally:

    $ go get github.com/sr/grpcinstrument/...
    $ godoc github.com/sr/grpcinstrument

To generate instrumentation stubs first install the compiler plugin:

    $ go install github.com/sr/grpcinstrument/...

Then, assuming your proto files are located under `src/`, invoke the `protoc`
command like so:

    $ protoc --grpcinstrument_out=src/ src/*.proto
    $ ls src/*-gen.go
    src/instrumented_buildkite-gen.go

Then update your main package to use the generated instrumented code:

```diff
--- a.go        2015-12-15 13:01:03.307634424 +0000
+++ b.go        2015-12-15 13:00:40.386861334 +0000
@@ -1,3 +1,5 @@
 server := grpc.NewServer()
+instrumentator := myapp.NewInstrumentator()
 buildkiteServer, := buildkite.NewAPIServer()
-buildkite.RegisterBuildkiteServiceServer(server, buildkiteServer)
+instrumented := buildkite.NewInstrumentedBuildkiteServiceServer(instrumentator, buildkiteServer)
+buildkite.RegisterBuildkiteServiceServer(server, instrumented)
```


package example

import (
	"time"
	"golang.org/x/net/context"
	"github.com/sr/grpcinstrument"
)


// InstrumentedExampleServer implements and instruments ExampleServer
// using the grpcinstrument package.
type InstrumentedExampleServer struct {
	instrumentator grpcinstrument.Instrumentator
	server ExampleServer
}


// NewInstrumentedExampleServer constructs a instrumentation wrapper for
// ExampleServer.
func NewInstrumentedExampleServer(
	instrumentator grpcinstrument.Instrumentator,
	server ExampleServer,
) *InstrumentedExampleServer {
	return &InstrumentedExampleServer{
		instrumentator,
		server,
	}
}

// Hello instruments the ExampleServer.Hello method.
func (a *InstrumentedExampleServer) Hello(
	ctx context.Context,
	request *HelloRequest,
) (response *HelloResponse, err error) {
	defer func(start time.Time) {
		grpcinstrument.Instrument(
			a.instrumentator,
			"example",
			"Hello",
			"HelloRequest",
			"HelloResponse",
			err,
			start,
		)
	}(time.Now())
	return a.server.Hello(ctx, request)
}

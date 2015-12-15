
package example

import (
	"time"
	"golang.org/x/net/context"
	"github.com/sr/grpcinstrument"
)


type instrumentedExampleServer struct {
	instrumentator grpcinstrument.Instrumentator
	server ExampleServer
}


func NewInstrumentedExampleServer(
	instrumentator grpcinstrument.Instrumentator,
	server ExampleServer,
) *instrumentedExampleServer {
	return &instrumentedExampleServer{
		instrumentator,
		server,
	}
}

func (a *instrumentedExampleServer) Hello(
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

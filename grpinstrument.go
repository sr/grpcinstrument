package grpcinstrument

import (
	"fmt"
	"time"

	"github.com/rcrowley/go-metrics"
	"go.pedge.io/proto/time"
)

type Logger interface {
	Log(*Call)
}

func Instrument(
	logger Logger,
	registry metrics.Registry,
	serviceName string,
	methodName string,
	inputType string,
	outputType string,
	err error,
	start time.Time,
) {
	call := &Call{
		Service:  serviceName,
		Method:   methodName,
		Input:    &Input{Type: inputType},
		Output:   &Output{Type: outputType},
		Duration: prototime.DurationToProto(time.Since(start)),
	}
	if err != nil {
		call.Error = &Error{Message: err.Error()}
	}
	logger.Log(call)
	metrics.GetOrRegisterCounter(fmt.Sprintf("%s.methods.%s.calls", call.Service, call.Method), registry).Inc(1)
}

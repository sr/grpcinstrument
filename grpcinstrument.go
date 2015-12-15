package grpcinstrument

import (
	"time"

	"go.pedge.io/proto/time"
)

type Instrumentator interface {
	Init() error
	Log(*Call)
	Measure(*Call)
}

func Instrument(
	instrumentator Instrumentator,
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
	instrumentator.Measure(call)
	instrumentator.Log(call)
}

func (c *Call) IsError() bool {
	return c.Error != nil
}

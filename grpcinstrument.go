package grpcinstrument

import (
	"fmt"
	"time"

	"go.pedge.io/proto/time"
)

type Instrumentator interface {
	Log(*Call)
	Increment(counterName string)
}

func (c *Call) IsError() bool {
	return c.Error != nil
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
	metric := fmt.Sprintf("%s.methods.%s", call.Service, call.Method)
	if call.IsError() {
		instrumentator.Increment(fmt.Sprintf("%s.errors", metric))
	}
	instrumentator.Increment(fmt.Sprintf("%s.calls", metric))
	instrumentator.Increment(fmt.Sprintf("%s.success", metric))
	instrumentator.Log(call)
}

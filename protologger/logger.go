package protologger

import (
	"github.com/sr/grpcinstrument"
	"go.pedge.io/protolog"
)

type logger struct {
	log protolog.Logger
}

func newLogger(log protolog.Logger) grpcinstrument.Logger {
	return &logger{log}
}

// Init does nothing. It exists only to satisfy the Logger interface.
func (l *logger) Init() error {
	return nil
}

// Log writes the RPC call to protolog using the INFO severity level for
// successful RPC calls or ERROR if the call resulted in an error.
func (l *logger) Log(call *grpcinstrument.Call) {
	if call.IsError() {
		l.log.Error(call)
	} else {
		l.log.Info(call)
	}
}

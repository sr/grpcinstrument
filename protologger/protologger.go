package protologger

import (
	"github.com/sr/grpcinstrument"
	"go.pedge.io/protolog"
)

// NewLogger constructs an implementation of the Logger interface that logs
// RPC calls via protolog.
func NewLogger(logger protolog.Logger) grpcinstrument.Logger {
	return newLogger(logger)
}

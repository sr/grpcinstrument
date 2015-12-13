package grpcinstrument

import (
	"encoding/json"
	"log"
)

type logLogger struct{}

func newLogLogger() *logLogger {
	return &logLogger{}
}

func (l *logLogger) Log(call *Call) {
	encoded, _ := json.Marshal(call)
	log.Print(string(encoded))
}

package generator

import "text/template"

type fileDescriptor struct {
	Package  string
	Services []*serviceDescriptor
}

type serviceDescriptor struct {
	ServerInterface string
	Methods         []*methodDescriptor
}

type methodDescriptor struct {
	Service         string
	ServerInterface string
	Name            string
	InputType       string
	OutputType      string
}

var loggerTemplate = template.Must(template.New("instrumented_api_server.go").Parse(`
package {{.Package}}

import (
	"github.com/sr/grpcinstrument"
	"golang.org/x/net/context"
	"time"
)

{{range .Services}}
// Instrumented{{.ServerInterface}} implements and instruments {{.ServerInterface}}
// using the grpcinstrument package.
type Instrumented{{.ServerInterface}} struct {
	instrumentator grpcinstrument.Instrumentator
	server         {{.ServerInterface}}
}
{{end}}
{{range .Services}}
// NewInstrumented{{.ServerInterface}} constructs a instrumentation wrapper for
// {{.ServerInterface}}.
func NewInstrumented{{.ServerInterface}}(
	instrumentator grpcinstrument.Instrumentator,
	server {{.ServerInterface}},
) *Instrumented{{.ServerInterface}} {
	return &Instrumented{{.ServerInterface}}{
		instrumentator,
		server,
	}
}
{{range .Methods}}
// {{.Name}} instruments the {{.ServerInterface}}.{{.Name}} method.
func (a *Instrumented{{.ServerInterface}}) {{.Name}}(
	ctx context.Context,
	request *{{.InputType}},
) (response *{{.OutputType}}, err error) {
	defer func(start time.Time) {
		grpcinstrument.Instrument(
			a.instrumentator,
			"{{.Service}}",
			"{{.Name}}",
			"{{.InputType}}",
			"{{.OutputType}}",
			err,
			start,
		)
	}(time.Now())
	return a.server.{{.Name}}(ctx, request)
}
{{end}}{{end}}`))

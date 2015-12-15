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
	"time"
	"golang.org/x/net/context"
	"github.com/sr/operator/src/grpcinstrument"
)

{{range .Services}}
type instrumented{{.ServerInterface}} struct {
	instrumentator grpcinstrument.Instrumentator
	server {{.ServerInterface}}
}
{{end}}
{{range .Services}}
func NewInstrumented{{.ServerInterface}}(
	instrumentator grpcinstrument.Instrumentator,
	server {{.ServerInterface}},
) *instrumented{{.ServerInterface}} {
	return &instrumented{{.ServerInterface}}{
		instrumentator,
		server,
	}
}
{{range .Methods}}
func (a *instrumented{{.ServerInterface}}) {{.Name}}(
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

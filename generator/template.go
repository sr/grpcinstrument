package generator

import "text/template"

type fileDescriptor struct {
	Package         string
	Service         string
	ServerInterface string
	Methods         []*methodDescriptor
}

type methodDescriptor struct {
	Service    string
	Name       string
	InputType  string
	OutputType string
}

var loggerTemplate = template.Must(template.New("instrumented_api_server.go").Parse(`
package {{.Package}}

import (
	"time"
	"golang.org/x/net/context"
	"github.com/sr/operator/src/grpcinstrument"
	"github.com/rcrowley/go-metrics"
)

type instrumentedAPIServer struct {
	instrumentator grpcinstrument.Instrumentator
	delegate {{.ServerInterface}}
}

func NewInstrumentedAPIServer(
	instrumentator grpcinstrument.Instrumentator,
	delegate {{.ServerInterface}},
) *instrumentedAPIServer {
	return &instrumentedAPIServer{
		instrumentator,
		delegate,
	}
}

{{range .Methods}}
func (a *instrumentedAPIServer) {{.Name}}(
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
	return a.delegate.{{.Name}}(ctx, request)
}
{{end}}`))

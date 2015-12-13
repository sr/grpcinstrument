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

var loggerTemplate = template.Must(template.New("log_api_server.go").Parse(`
package {{.Package}}

import (
	"time"
	"golang.org/x/net/context"
	"github.com/sr/operator/src/grpclog"
	"github.com/rcrowley/go-metrics"
)

type logAPIServer struct {
	logger grpclog.Logger
	metrics metrics.Registry
	delegate {{.ServerInterface}}
}

func NewLogAPIServer(
	logger grpclog.Logger,
	metrics metrics.Registry,
	delegate {{.ServerInterface}},
) *logAPIServer {
	return &logAPIServer{logger, metrics, delegate}
}

{{range .Methods}}
func (a *logAPIServer) {{.Name}}(
	ctx context.Context,
	request *{{.InputType}},
) (response *{{.OutputType}}, err error) {
	defer func(start time.Time) {
		grpclog.Instrument(
			a.logger,
			a.metrics,
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

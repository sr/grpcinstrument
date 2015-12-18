package example

import "golang.org/x/net/context"

// ExampleServer is an example gRPC server.
type ExampleServer struct{}

// Hello returns the name of the author's favourite DJ.
func (a *ExampleServer) Hello(
	ctx context.Context,
	request *HelloRequest,
) (response *HelloResponse, err error) {
	return &HelloResponse{Message: "Tini Günter"}, nil
}

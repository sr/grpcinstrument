package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gogo/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/sr/operator/src/grpcinstrument/generator"
)

func fatal(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Print("protoc-gen-grpclog error:", s)
	os.Exit(1)
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatal(err, "reading input")
	}
	request := &plugin.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, request); err != nil {
		fatal(err, "parsing input proto")
	}
	gen := generator.NewGenerator(request)
	response, err := gen.Generate()
	if err != nil {
		fatal(err, "generating response")
	}
	data, err = proto.Marshal(response)
	if err != nil {
		fatal(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		fatal(err, "failed to write output proto")
	}
}

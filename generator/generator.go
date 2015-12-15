package generator

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

const (
	prefix = "instrumented_"
	suffix = "-gen.go"
)

type Generator struct {
	request *plugin.CodeGeneratorRequest
}

func NewGenerator(request *plugin.CodeGeneratorRequest) *Generator {
	return &Generator{request}
}

func (g *Generator) Generate() (*plugin.CodeGeneratorResponse, error) {
	numFiles := len(g.request.FileToGenerate)
	if numFiles == 0 {
		return nil, errors.New("no file to generate")
	}
	response := &plugin.CodeGeneratorResponse{}
	response.File = make([]*plugin.CodeGeneratorResponse_File, numFiles)
	filesByName := make(map[string]*descriptor.FileDescriptorProto, numFiles)
	for _, file := range g.request.ProtoFile {
		filesByName[file.GetName()] = file
	}
	for i, filePath := range g.request.FileToGenerate {
		content, err := g.generateFile(filesByName[filePath])
		if err != nil {
			return nil, err
		}
		dir, file := path.Split(filePath)
		genFile := strings.Replace(file, path.Ext(file), suffix, 1)
		response.File[i] = &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("%s/%s%s", dir, prefix, genFile)),
			Content: proto.String(content),
		}
	}
	return response, nil
}

func (g *Generator) generateFile(
	file *descriptor.FileDescriptorProto,
) (content string, err error) {
	descriptor := &fileDescriptor{
		Package:  file.GetPackage(),
		Services: make([]*serviceDescriptor, len(file.Service)),
	}
	for i, service := range file.Service {
		serviceDescriptor := &serviceDescriptor{
			ServerInterface: fmt.Sprintf("%sServer", service.GetName()),
			Methods:         make([]*methodDescriptor, len(service.Method)),
		}
		for n, method := range service.Method {
			serviceDescriptor.Methods[n] = &methodDescriptor{
				Service:         descriptor.Package,
				ServerInterface: serviceDescriptor.ServerInterface,
				Name:            method.GetName(),
				InputType:       strings.Split(method.GetInputType(), ".")[2],
				OutputType:      strings.Split(method.GetOutputType(), ".")[2],
			}
		}
		descriptor.Services[i] = serviceDescriptor
	}
	var buffer bytes.Buffer
	if err := loggerTemplate.Execute(&buffer, descriptor); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

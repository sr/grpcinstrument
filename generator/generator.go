package generator

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/gogo/protobuf/proto"
	google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	plugin "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
)

type Generator interface {
	Generate() (*plugin.CodeGeneratorResponse, error)
}

func NewGenerator(request *plugin.CodeGeneratorRequest) Generator {
	return newGenerator(request)
}

type generator struct {
	request        *plugin.CodeGeneratorRequest
	messagesByName map[string]*google_protobuf.DescriptorProto
}

func newGenerator(request *plugin.CodeGeneratorRequest) *generator {
	return &generator{
		request,
		make(map[string]*google_protobuf.DescriptorProto),
	}
}

func (g *generator) Generate() (*plugin.CodeGeneratorResponse, error) {
	numFiles := len(g.request.FileToGenerate)
	if numFiles == 0 {
		return nil, errors.New("no file to generate")
	}
	response := &plugin.CodeGeneratorResponse{}
	response.File = make([]*plugin.CodeGeneratorResponse_File, numFiles)
	filesByName := make(map[string]*google_protobuf.FileDescriptorProto, numFiles)
	for _, file := range g.request.ProtoFile {
		filesByName[file.GetName()] = file
	}
	for i, filePath := range g.request.FileToGenerate {
		content, err := g.generateFile(filesByName[filePath])
		if err != nil {
			return nil, err
		}
		dir, file := path.Split(filePath)
		genFile := strings.Replace(file, path.Ext(file), "-gen.go", 1)
		response.File[i] = &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("%s/log_%s", dir, genFile)),
			Content: proto.String(content),
		}
	}
	return response, nil
}

func (g *generator) generateFile(
	file *google_protobuf.FileDescriptorProto,
) (content string, err error) {
	if len(file.Service) != 1 {
		return "", errors.New("can only generate script for exactly one service")
	}
	service := file.Service[0]
	descriptor := &fileDescriptor{
		Service:         file.GetPackage(),
		Package:         file.GetPackage(),
		ServerInterface: fmt.Sprintf("%sServer", service.GetName()),
		Methods:         make([]*methodDescriptor, len(service.Method)),
	}
	for i, method := range service.Method {
		descriptor.Methods[i] = newMethodDescriptor(file.GetPackage(), method)
	}
	var buffer bytes.Buffer
	if err := loggerTemplate.Execute(&buffer, descriptor); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func newMethodDescriptor(
	service string,
	method *google_protobuf.MethodDescriptorProto,
) *methodDescriptor {
	return &methodDescriptor{
		Service:    service,
		Name:       method.GetName(),
		InputType:  strings.Split(method.GetInputType(), ".")[2],
		OutputType: strings.Split(method.GetOutputType(), ".")[2],
	}
}

proto:
	go get github.com/sr/grpcinstrument/...
	go get -v go.pedge.io/protoeasy/cmd/protoeasy
	go get -v go.pedge.io/pkg/cmd/strip-package-comments
	protoeasy --go --grpc --go-import-path github.com/sr/grpcinstrument/example .
	find . -name *\.pb\*\.go | xargs strip-package-comments
	protoc --grpcinstrument_out=. *.proto

.PHONY: \
	proto

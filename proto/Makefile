# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
TARGETS=eventbroker servicemesh
PKGPATH=github.com/dan-solli/homeapps/proto

# Targets
all: protobuf

.PHONY: protobuf
protobuf: $(TARGETS)

eventbroker: eventbroker.proto
	@mkdir -p ./$@
	protoc --go_out=./$@ --go_opt=paths=source_relative --go-grpc_out=./$@ --go-grpc_opt=paths=source_relative $@.proto
	@(cd $@ && go mod init $(PKGPATH)/$@ && go mod tidy)

servicemesh: servicemesh.proto
	@mkdir -p ./$@
	protoc --go_out=./$@ --go_opt=paths=source_relative --go-grpc_out=./$@ --go-grpc_opt=paths=source_relative $@.proto
	@(cd $@ && go mod init $(PKGPATH)/$@ && go mod tidy)

clean:
	rm -rf $(TARGETS)
	rm -f *.go

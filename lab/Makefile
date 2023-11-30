# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
TEST_FLAGS=-race -coverprofile=coverage.out -covermode=atomic

# Targets
all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) $(TEST_FLAGS) ./...

clean:
	rm -f $(BINARY_NAME) coverage.out

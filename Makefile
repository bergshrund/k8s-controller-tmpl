BINARY_NAME := k8s-controller
VERSION := $(shell git describe --tags --always --dirty)
BUILD_FLAGS = -v -o $(BINARY_NAME) -ldflags "-X=k8s-controller-tmpl/cmd.appVersion=$(VERSION)"
TARGETOS = "linux"
TARGETARCH = "amd64"
REGISTRY = bergshrund

.PHONY: all build run clean envtest

all: build

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
ENVTEST_VERSION ?= latest

format:
	gofmt -s -w ./

lint:
	go vet ./...

envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS)

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t ${REGISTRY}/${BINARY_NAME}:${VERSION}-${TARGETOS}-${TARGETARCH} .

run:
	go run main.go

clean:
	rm -f $(BINARY_NAME)

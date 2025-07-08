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

# MSYS_NO_PATHCONV=1 is used to prevent MSYS from converting paths, which is necessary for Windows compatibility.
test: envtest
	go install gotest.tools/gotestsum@latest
	KUBEBUILDER_ASSETS="$(shell MSYS_NO_PATHCONV=1 $(ENVTEST) use --bin-dir $(LOCALBIN) -p path)" gotestsum --junitfile report.xml --format testname ./... ${TEST_ARGS}

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t ${REGISTRY}/${BINARY_NAME}:${VERSION}-${TARGETOS}-${TARGETARCH} .

run:
	go run main.go

clean:
	rm -f $(BINARY_NAME)

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef
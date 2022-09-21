GO := go

GO_BUILD_BINDIR :=./bin
GIT_COMMIT := $(or $(SOURCE_GIT_COMMIT),$(shell git rev-parse --short HEAD))
GIT_TAG :="$(shell git tag | sort -V | tail -1)"

GO_LD_EXTRAFLAGS :=-X github.com/uor-community/argocd-uor-plugin/cmd.version="$(shell git tag | sort -V | tail -1)" \
				   -X github.com/uor-community/argocd-uor-plugin/cmd.buildData="dev" \
				   -X github.com/uor-community/argocd-uor-plugin/cmd.commit="$(GIT_COMMIT)" \
				   -X github.com/uor-community/argocd-uor-plugin/cmd.buildDate="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"

build: prep-build-dir
	$(GO) build -o $(GO_BUILD_BINDIR)/argocd-uor-plugin  -ldflags="$(GO_LD_EXTRAFLAGS)" main.go
.PHONY: build

prep-build-dir:
	mkdir -p ${GO_BUILD_BINDIR}
.PHONY: prep-build-dir

vendor:
	$(GO) mod tidy
	$(GO) mod verify
	$(GO) mod vendor
.PHONY: vendor

clean:
	@rm -rf ./$(GO_BUILD_BINDIR)/*
.PHONY: clean

sanity: vendor format vet
	git diff --exit-code
.PHONY: sanity

format: 
	$(GO) fmt ./...
.PHONY: format

vet: 
	$(GO) vet ./...
.PHONY: vet

all: clean vendor build
.PHONY: all

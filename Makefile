# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/sapcc/digicert-issuer
CRD_OPTIONS ?= "crd:crdVersions=v1"

# Temporary directory for tools
TOOLS_BIN_DIR = $(shell pwd)/tmp/bin

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: generate manifests build

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out | grep -v "no test files"

# Build digicert-issuer binary.
build: BUILD_DATE  = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
build: GIT_REVISION  = $(shell git rev-parse --short HEAD)
build: GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
build: VERSION=$(shell cat VERSION)
build: fmt vet
	go build\
		-ldflags "-s -w -X github.com/sapcc/digicert-issuer/pkg/version.Revision=$(GIT_REVISION) -X github.com/sapcc/digicert-issuer/pkg/version.Branch=$(GIT_BRANCH) -X github.com/sapcc/digicert-issuer/pkg/version.BuildDate=$(BUILD_DATE) -X github.com/sapcc/digicert-issuer/pkg/version.Version=$(VERSION)"\
		-o bin/digicert-issuer main.go

# Make sure to run the deploy-local target beforehand.
# Run the digicert-issuer against the configured Kubernetes cluster in ~/.kube/config with debug logging enabled.
run: generate build
	./bin/digicert-issuer

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy everything except the digicert-issuer to the configured Kubernetes cluster.
# The digicert-issuer itself is being run locally.
deploy-local: VERSION=$(shell cat VERSION)
deploy-local: manifests
	cd config/digicert-issuer && kustomize edit set image digicert-issuer=${IMG}:${VERSION}
	kustomize build config/default-local-operator | kubectl apply -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: VERSION=$(shell cat VERSION)
deploy: manifests
	cd config/digicert-issuer && kustomize edit set image digicert-issuer=${IMG}:${VERSION}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: _manifests apidocs

_manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=digicert-issuer-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: VERSION=$(shell cat VERSION)
docker-build:
	docker build . -t ${IMG}:${VERSION}

# Push the docker image
docker-push: VERSION=$(shell cat VERSION)
docker-push:
	docker push ${IMG}:${VERSION} &&\
	docker tag ${IMG}:${VERSION} ${IMG}:latest &&\
	docker push ${IMG}:latest

apidocs: doc-gen
	$(DOC_GEN) apis/certmanager/v1beta1/*.go > docs/apidocs/api.md

git-push-tag: VERSION=$(shell cat VERSION)
git-push-tag:
	git push origin ${VERSION}

git-tag-release: VERSION=$(shell cat VERSION)
git-tag-release: check-release-version
	git tag --annotate ${VERSION} --message "digicert-issuer ${VERSION}"

check-release-version: VERSION=$(shell cat VERSION)
check-release-version:
	if test x$$(git tag --list ${VERSION}) != x; \
	then \
		echo "Tag [${VERSION}] already exists. Please check the working copy."; git diff . ; exit 1;\
	fi

set-image: VERSION=$(shell cat VERSION)
set-image:
	cd config/digicert-issuer && kustomize edit set image digicert-issuer=${IMG}:${VERSION}
	git commit --allow-empty -am "set digicert-issuer image to ${VERSION}"

release: VERSION=$(shell cat VERSION)
release: git-tag-release set-image git-push-tag docker-build docker-push

$(TOOLS_BIN_DIR):
	mkdir -p $(TOOLS_BIN_DIR)

clean:
	rm -rf $(TOOLS_BIN_DIR) ./bin

.PHONY: vendor
vendor:
	go mod download

# Find or download controller-gen.
controller-gen: VERSION=v0.3.0
controller-gen: $(TOOLS_BIN_DIR)
ifeq (,$(wildcard $(TOOLS_BIN_DIR)/controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	GOBIN=$(TOOLS_BIN_DIR) go get sigs.k8s.io/controller-tools/cmd/controller-gen@${VERSION} ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
endif
CONTROLLER_GEN=$(TOOLS_BIN_DIR)/controller-gen

doc-gen: $(TOOLS_BIN_DIR)
ifeq (,$(wildcard $(TOOLS_BIN_DIR)/doc-gen))
	@{ \
	set -e ;\
	GOBIN=$(TOOLS_BIN_DIR) go install ./cmd/doc-gen ;\
	}
endif
DOC_GEN=$(TOOLS_BIN_DIR)/doc-gen

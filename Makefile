# SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
#
# SPDX-License-Identifier: Apache-2.0

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

build-all: generate manifests build

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
		-o bin/digicert-issuer cmd/digicert-issuer/main.go

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

.PHONY: vendor
vendor:
	go mod download

## Tool Binaries
CONTROLLER_TOOLS_VERSION ?= v0.17.0
CONTROLLER_GEN ?= $(TOOLS_BIN_DIR)/controller-gen
DOC_GEN ?= $(TOOLS_BIN_DIR)/doc-gen

$(TOOLS_BIN_DIR):
	mkdir -p $(TOOLS_BIN_DIR)

clean:
	rm -rf $(TOOLS_BIN_DIR) ./bin

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(TOOLS_BIN_DIR)
	GOBIN=$(TOOLS_BIN_DIR) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: doc-gen
doc-gen: $(DOC_GEN)
$(DOC_GEN): $(TOOLS_BIN_DIR)
	GOBIN=$(TOOLS_BIN_DIR) go install ./cmd/doc-gen

UNAME_S := $(shell uname -s)
SED = sed
XARGS = xargs
ifeq ($(UNAME_S),Darwin)
	SED = gsed
	XARGS = gxargs
endif

GO_COVERPKGS := $(shell go list ./...)
GO_TESTPKGS := ./...

tidy-deps: FORCE
	go mod tidy
	go mod verify

build/cover.out: FORCE install-ginkgo generate install-setup-envtest | build
	@printf "\e[1;36m>> Running tests\e[0m\n"
	KUBEBUILDER_ASSETS=$$(setup-envtest use 1.32 -p path) ginkgo run --randomize-all -output-dir=build $(GO_BUILDFLAGS) -ldflags '-s -w $(GO_LDFLAGS)' -covermode=count -coverpkg=$(subst $(space),$(comma),$(GO_COVERPKGS)) $(GO_TESTPKGS)
	@mv build/coverprofile.out build/cover.out

build/cover.html: build/cover.out
	@printf "\e[1;36m>> go tool cover > build/cover.html\e[0m\n"
	@go tool cover -html $< -o $@

install-ginkgo: FORCE
	@if ! hash ginkgo 2>/dev/null; then printf "\e[1;36m>> Installing ginkgo (this may take a while)...\e[0m\n"; go install github.com/onsi/ginkgo/v2/ginkgo@latest; fi

install-setup-envtest: FORCE
	@if ! hash setup-envtest 2>/dev/null; then printf "\e[1;36m>> Installing setup-envtest (this may take a while)...\e[0m\n"; go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest; fi

check-dependency-licenses: FORCE install-go-licence-detector
	@printf "\e[1;36m>> go-licence-detector\e[0m\n"
	@go list -m -mod=readonly -json all | go-licence-detector -includeIndirect -rules .license-scan-rules.json -overrides .license-scan-overrides.jsonl

install-go-licence-detector: FORCE
	@if ! hash go-licence-detector 2>/dev/null; then printf "\e[1;36m>> Installing go-licence-detector (this may take a while)...\e[0m\n"; go install go.elastic.co/go-licence-detector@latest; fi

install-addlicense: FORCE
	@if ! hash addlicense 2>/dev/null; then printf "\e[1;36m>> Installing addlicense (this may take a while)...\e[0m\n"; go install github.com/google/addlicense@latest; fi

license-headers: FORCE install-addlicense
	@printf "\e[1;36m>> addlicense (for license headers on source code files)\e[0m\n"
	@printf "%s\0" $(patsubst $(shell awk '$$1 == "module" {print $$2}' go.mod)%,.%/*.go,$(shell go list ./...)) | $(XARGS) -0 -I{} bash -c 'year="$$(grep 'Copyright' {} | head -n1 | grep -E -o '"'"'[0-9]{4}(-[0-9]{4})?'"'"')"; gawk -i inplace '"'"'{if (display) {print} else {!/^\/\*/ && !/^\*/}}; {if (!display && $$0 ~ /^(package |$$)/) {display=1} else { }}'"'"' {}; addlicense -c "SAP SE or an SAP affiliate company" -s=only -y "$$year" -- {}; $(SED) -i '"'"'1s+// Copyright +// SPDX-FileCopyrightText: +'"'"' {}'
	@printf "\e[1;36m>> reuse annotate (for license headers on other files)\e[0m\n"
	@reuse lint -j | jq -r '.non_compliant.missing_licensing_info[]' | grep -vw vendor | $(XARGS) reuse annotate -c 'SAP SE or an SAP affiliate company' -l Apache-2.0 --skip-unrecognised
	@printf "\e[1;36m>> reuse download --all\e[0m\n"
	@reuse download --all
	@printf "\e[1;35mPlease review the changes. If *.license files were generated, consider instructing go-makefile-maker to add overrides to REUSE.toml instead.\e[0m\n"

check-license-headers: FORCE install-addlicense tidy-deps
	@printf "\e[1;36m>> addlicense --check\e[0m\n"
	@addlicense --check -- $(patsubst $(shell awk '$$1 == "module" {print $$2}' go.mod)%,.%/*.go,$(shell go list ./...))

.PHONY: FORCE

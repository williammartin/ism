
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
SM = bin/sm

all: clean test manager cli

# Run tests
test: generate fmt vet manifests kubebuilder-tests acceptance-tests

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager github.com/pivotal-cf/ism/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
	go generate ./pkg/... ./cmd/...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

### CUSTOM MAKE RULES ###

# Make the CLI
cli:
	go build -o ${SM} cmd/sm/main.go

# Clean the CLI
clean:
	rm -f ${SM}

# Run acceptance tests
acceptance-tests:
	ginkgo -r acceptance

# Run kubebuilder tests
kubebuilder-tests:
	go test ./pkg/... ./cmd/... -coverprofile cover.out

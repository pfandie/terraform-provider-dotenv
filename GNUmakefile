.PHONY: default fmt lint build install generate test testacc testacc-tofu

default: fmt lint install generate

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

build:
	go build -v ./...

install: build
	go install -v ./...

# generates the docs; tfplugindocs is pinned via the go.mod tool directive
generate:
	go tool tfplugindocs generate

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./internal/provider/

testacc-tofu:
	TF_ACC=1 TF_ACC_TERRAFORM_PATH=$$(which tofu) TF_ACC_PROVIDER_HOST=registry.opentofu.org \
		go test -v -cover -timeout 120m ./internal/provider/

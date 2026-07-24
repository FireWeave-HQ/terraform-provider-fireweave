default: test

build:
	go build -o terraform-provider-fireweave .

test:
	go test ./...

testacc:
	TF_ACC=1 go test ./internal/provider -v -timeout 30m

docs:
	go generate ./...

.PHONY: build test testacc docs

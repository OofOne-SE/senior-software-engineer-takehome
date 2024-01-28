.PHONY: download generate

SWAGGER=v0.24.0

download:
	- echo Download go.mod dependencies
	- go mod download

generate:
	docker run --rm -it --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v ${PWD}:/go quay.io/goswagger/swagger:${SWAGGER} generate server -f spec/openapi.yaml
	docker run --rm -it --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v ${PWD}:/go quay.io/goswagger/swagger:${SWAGGER} generate client -f spec/openapi.yaml

validate: 
	docker run --rm -it --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v ${PWD}:/go quay.io/goswagger/swagger:${SWAGGER} validate spec/openapi.yaml

fix: 
	golangci-lint run --fix

showspec:
	docker run --rm -it --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v ${PWD}:/go quay.io/goswagger/swagger:${SWAGGER} serve -p 8088 ./spec/openapi.yaml

test:
	go clean -testcache && go test ./...
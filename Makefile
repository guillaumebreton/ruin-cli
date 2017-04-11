build:
	go build .

fmt:
	go fmt . && go vet .

install: build
	go install

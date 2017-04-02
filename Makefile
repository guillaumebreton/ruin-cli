build:
	go build .

fmt:
	go fmt . && go vet .

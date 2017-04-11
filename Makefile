build:
	go build .

fmt:
	go fmt . && go vet .

install:
	go install

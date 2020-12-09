.PHONY: setup test

TAG=latest

setup:
	go get -v ./...
	go get github.com/golang/lint/golint
	go get github.com/fzipp/gocyclo

test:
	go test -v ./...
	go test -cover ./...
	golint ./...
	go tool vet -all .
	gocyclo -over 10 .

clean:
	rm kacl-*.tar.gz

build:
	go build -o kacl main.go

releases:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o kacl main.go    && tar czvf kacl-$(TAG)-linux_amd64.tar.gz kacl
	GOARCH=386   GOOS=linux CGO_ENABLED=0 go build -o kacl main.go    && tar czvf kacl-$(TAG)-linux_386.tar.gz kacl
	GOARCH=amd64 GOOS=darwin CGO_ENABLED=0 go build -o kacl main.go   && tar czvf kacl-$(TAG)-darwin_amd64.tar.gz kacl
	GOARCH=386   GOOS=darwin CGO_ENABLED=0 go build -o kacl main.go   && tar czvf kacl-$(TAG)-darwin_386.tar.gz kacl
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -o kacl main.go  && tar czvf kacl-$(TAG)-windows_amd64.tar.gz kacl
	GOARCH=386   GOOS=windows CGO_ENABLED=0 go build -o kacl main.go  && tar czvf kacl-$(TAG)-windows_386.tar.gz kacl
	rm ./kacl

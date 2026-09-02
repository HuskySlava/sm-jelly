BINARY := sm-jelly

.PHONY: build run fmt vet tidy check clean

build:
	go build -o bin/$(BINARY) .

run:
	go run .

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

check: fmt vet build

clean:
	rm -rf bin
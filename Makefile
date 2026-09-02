BINARY := sm-jelly

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
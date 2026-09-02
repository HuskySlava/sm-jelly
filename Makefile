BINARY := sm-jelly
IMAGE    := sm-jelly
PLATFORM := linux/amd64

.PHONY: build run fmt vet tidy check clean docker-build

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

docker-build:
	docker build --platform $(PLATFORM) -t $(IMAGE) .
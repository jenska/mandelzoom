APP := mandelzoom

.PHONY: all run build tidy fmt test clean

all: build

run:
	go run ./...

build:
	go build -o bin/$(APP) ./...

tidy:
	go mod tidy

fmt:
	gofmt -w *.go

test:
	go test ./...

clean:
	rm -rf bin

TEMPL_BIN ?= templ
TAILWIND_BIN ?= npx @tailwindcss/cli
BINARY_NAME ?= bin/orbitodx

.PHONY: all build build-vps generate css run test clean

all: build

generate:
	$(TEMPL_BIN) generate -path ./views

css:
	$(TAILWIND_BIN) -i ./web/input.css -o ./web/static/css/app.css --minify

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/orbitodx

build-vps:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/orbitodx

run: build
	./$(BINARY_NAME) --host 0.0.0.0 --port 8080

test:
	go test -v ./internal/...

clean:
	rm -rf bin/

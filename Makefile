MODULE := $(shell grep '^module ' go.mod | cut -d' ' -f2)

.PHONY: help install-tools
help:
	@grep -E '^[a-zA-Z][a-zA-Z0-9_.-]*:([^=]|$$)' $(MAKEFILE_LIST) | cut -d: -f1 | sort -u

install-tools:
	go install github.com/a-h/templ/cmd/templ@v0.2.793
	go install golang.org/x/tools/cmd/goimports@v0.30.0
	go install github.com/air-verse/air@v1.65.3

start:
	docker compose up -d --build

stop:
	docker compose down --remove-orphans

reload: stop start

logs:
	docker compose logs -f

test:
	go test -v -race -count=1 -timeout=60s -shuffle=on -failfast -cover ./...

goimports:
	goimports -w -local $(MODULE) .

templ:
	templ fmt .
	templ generate

build-local: goimports templ build-wasm
	go build -o ./tmp/main ./cmd/image-text 

build-wasm:
	GOOS=js GOARCH=wasm go build -o wasm/main.wasm ./cmd/wasm/main.go

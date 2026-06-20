MODULE := $(shell grep '^module ' go.mod | cut -d' ' -f2)
COMMIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
LDFLAGS := -X $(MODULE)/version.Commit=$(COMMIT_SHA)

.PHONY: help install-tools
help:
	@grep -E '^[a-zA-Z][a-zA-Z0-9_.-]*:([^=]|$$)' $(MAKEFILE_LIST) | cut -d: -f1 | sort -u

install-tools:
	go install github.com/a-h/templ/cmd/templ@v0.2.793
	go install golang.org/x/tools/cmd/goimports@v0.30.0
	go install github.com/air-verse/air@v1.65.3
	@command -v wasm-opt >/dev/null 2>&1 || echo "Install binaryen for wasm-opt (e.g. brew install binaryen)"

start:
	COMMIT_SHA=$(COMMIT_SHA) docker compose up -d --build

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
	go build -ldflags "$(LDFLAGS)" -o ./tmp/main ./cmd/image-text

build-wasm:
	GOOS=js GOARCH=wasm go build -o wasm/main.wasm ./cmd/wasm/main.go
	@if command -v wasm-opt >/dev/null 2>&1; then \
		wasm-opt -Oz --enable-bulk-memory --enable-nontrapping-float-to-int wasm/main.wasm -o wasm/main.wasm.tmp && \
		mv wasm/main.wasm.tmp wasm/main.wasm; \
	else \
		echo "wasm-opt not found; install binaryen for a smaller bundle (e.g. brew install binaryen)" >&2; \
	fi

start:
	docker compose up -d --build

stop:
	docker compose down

reload: stop start

logs:
	docker compose logs -f

test:
	go test ./...

goimports:
	goimports -w -local $(cat go.mod | grep "^module " | cut -d' ' -f 2) .

templ:
	templ fmt .
	templ generate

build-local: templ goimports build-wasm
	go build -o ./tmp/main ./cmd/image-text 

build-wasm:
	GOOS=js GOARCH=wasm go build -o wasm/main.wasm ./cmd/wasm/main.go

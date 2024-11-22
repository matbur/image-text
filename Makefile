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
	templ generate

build-local: templ goimports
	go build -o ./tmp/main ./cmd/image-text 

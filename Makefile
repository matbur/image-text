start:
	docker compose up -d --build

stop:
	docker compose down

reload: stop start

logs:
	docker compose logs -f

test:
	go test ./...
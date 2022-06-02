build:
	docker-compose up --build api

start:
	docker-compose up -d

migrate-init:
	docker run -v ./schema:/migrations --network host migrate/migrate     -path=/migrations/ -database postgres://postgres:qwerty123@localhost:54320/postgres?sslmode=disable up
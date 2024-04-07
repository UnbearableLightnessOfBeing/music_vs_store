postgres:
	docker run --name music_vs_store_postgres -p 1234:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:12-alpine

createdb:
	docker exec -it music_vs_store_postgres createdb --username=admin --owner=admin music_vs_store_db

dropdb:
	docker exec -it music_vs_store_postgres dropdb music_vs_store_db

migrate-up:
	migrate -path db/migrations -database "postgresql://admin:admin@localhost:1234/music_vs_store_db?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgresql://admin:admin@localhost:1234/music_vs_store_db?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test:
	go test ./... -v -cover

.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc test

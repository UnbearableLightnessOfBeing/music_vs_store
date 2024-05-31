include .env

#postgres:
	#docker run --name music_vs_store_postgres -p 1234:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:12-alpine

createdb:
	docker exec -it music_vs_store_db createdb --username=admin --owner=admin postgres_db

dropdb:
	docker exec -it music_vs_store_db dropdb postgres_db

migrate-up:
	migrate -path db/migrations -database "postgresql://admin:admin@localhost:1234/postgres_db?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgresql://admin:admin@localhost:1234/postgres_db?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test:
	go test ./... -v -cover

run-css-server:
	npx tailwindcss -i ./static/styles/input.css -o ./static/styles/output.css --watch

clear-images:
	rm -rf ./storage/images/*

copy-dump: 
	docker cp ./db/seeders/initial_seed.sql music_vs_store_db:/root/dump.sql

run-seeder:
	docker exec -it music_vs_store_db psql -U admin -d postgres_db -c '\i root/dump.sql'

reset-db:
	make migrate-down && make migrate-up && make clear-images && make copy-dump && make run-seeder

.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc test run-css-server clear-images copy-dump run-seeder reset-db

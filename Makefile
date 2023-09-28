DB_URL1=postgresql://root:secret@localhost:5432/clothing?sslmode=disable
SQLC_URL=D:\Study\EC_Clothing:/src
DB_URL=postgres://root:9OKjOaDKdKByUn8EzpAyjtzr4FPqH2hQ@dpg-ckafm3kg66mc73d3ul70-a.oregon-postgres.render.com/clothing_y3kn
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root clothing

dropdb:
	docker exec -it postgres12 dropdb clothing

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

gen:
	docker run --rm -v "$(SQLC_URL)" -w /src kjconroy/sqlc generate

server:
	go run main.go

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown gen new_migration test
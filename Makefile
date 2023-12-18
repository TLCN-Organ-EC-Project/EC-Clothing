DB_URL=postgresql://root:secret@localhost:5432/clothing?sslmode=disable
SQLC_URL=D:\Study\EC_Clothing:/src
DB_URL1=postgres://root:K68CtX7IIlOx5frordPam3mBIvroMfWu@dpg-clvu8t5a73kc73btbqf0-a.oregon-postgres.render.com/clothing_lbxj
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

migrateupR:
	migrate -path db/migration -database "$(DB_URL1)" -verbose up

migratedownR:
	migrate -path db/migration -database "$(DB_URL1)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migrateup1R:
	migrate -path db/migration -database "$(DB_URL1)" -verbose up 1

migratedown1R:
	migrate -path db/migration -database "$(DB_URL1)" -verbose down 1

gen:
	docker run --rm -v "$(SQLC_URL)" -w /src kjconroy/sqlc generate

server:
	go run main.go

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 migrateupR migratedownR migrateup1R migratedown1R gen new_migration test
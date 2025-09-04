createdb:
	docker exec -it learn-go-db-1 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it learn-go-db-1 dropdb --username=postgres simple_bank

migrateup:
	migrate --path db/migration --database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" --verbose up
	
migratedown:
	migrate --path db/migration --database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" --verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen --package mockdb --destination db/mock/store.go gilanggsb/simplebank/db/sqlc Store

.PHONY: createdb dropdb migrateup migratedown sqlc test server mock
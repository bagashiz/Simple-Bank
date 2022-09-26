postgres:
	docker run -h postgres-server --name postgres-server -p 5432:5432 --env-file ~/Docker/.postgres_env_list -d postgres
	
createdb:
	docker exec -ti postgres-server createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -ti postgres-server dropdb simple_bank
	
migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up
	
migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down
	
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
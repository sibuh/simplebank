postgres:
	docker run --name postgres21 --network bank-network -p 5432:5432 -e POSTGRES_USER=sibhatdb  -e POSTGRES_PASSWORD=sibhat21 -d postgres:12-alpine
createdb:
	docker exec -it postgres21 createdb --username=sibhatdb --owner=sibhatdb simple_bank
dropdb:
	docker exec -it postgres21 dropdb --username=sibhatdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://sibhatdb:sibhat21@localhost:5432/simple_bank?sslmode=disable" up 
migrateup1:
	migrate -path db/migration -database "postgresql://sibhatdb:sibhat21@localhost:5432/simple_bank?sslmode=disable" up 1


migratedown:
	migrate -path db/migration -database "postgresql://sibhatdb:sibhat21@localhost:5432/simple_bank?sslmode=disable" down 
migratedown1:
	migrate -path db/migration -database "postgresql://sibhatdb:sibhat21@localhost:5432/simple_bank?sslmode=disable" down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	 mockgen -package mockdb -destination db/mock/store.go assignment_01/simplebank/db/sqlc Store
.PHONY:postgres createdb dropdb migrateup migratedown sqlc test mock
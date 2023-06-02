postgres:
	docker run --name post -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -d postgres
createdb:
	docker exec -it post createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it post dropdb simple_bank
migrateup:
	migrate --path db/migrations -database "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate --path db/migrations -database "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrateup1:
	migrate --path db/migrations -database "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown1:
	migrate --path db/migrations -database "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
test:
	go test -v -cover ./...
cleantest:
	go clean -testcache
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/MikoBerries/SimpleBank/db/sqlc Store

# sqlcdocker:
# 	docker run --rm -v "C:\Users\Gio\Documents\goworkspace\src\github.com\MikoBerries\SimpleBank":/src -w /src kjconroy/sqlc generate
# sqlcdockerver:
# 	docker run --rm -v $(PWD):/src -w //src kjconroy/sqlc version
# sqlcdockerinit:
# 	docker run --rm --volumes-from myapps -v $(PWD):/src -w /src kjconroy/sqlc init
# sqlcdockergen:
# 	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc generate

PHONY: postgres createdb dropdb migrateup migratedown test cleantest server mock migrateup1 migratedown1

DB_URL=postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -d postgres
redis:
	docker run --name redis7 -p 6379:6379 -d redis:7.0-alpine --save 60 1 --loglevel warning
sb:
	docker run --name sb -p 8080:8080 simplebank:lastest
simpleBankRelease:
	docker run --name sb -p 8080:8080 -e GIN_MODE=release simplebank:lastest

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres15 dropdb simple_bank
migrateup:
	migrate --path db/migrations -database "$(DB_URL)" -verbose up
migratedown:
	migrate --path db/migrations -database "$(DB_URL) -verbose down
migrateup1:
	migrate --path db/migrations -database "$(DB_URL) -verbose up 1
migratedown1:
	migrate --path db/migrations -database "$(DB_URL) -verbose down 1

dbdocs:
	dbdocs build doc/db.dbml
dbml2sql:
	dbml2sql doc/db.dbml --postgres -o doc/schema.sql
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/MikoBerries/SimpleBank/db/sqlc Store
dockerBuild :
	docker build -t simplebank:lastest .

protoc:
	rm -f pb/*.go
	rm -f doc/swagger/*.json
	protoc --proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto 
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost -p 9090 -r repl


test:
	go test -v -cover ./...
cleantest:
	go clean -testcache
server:
	go run main.go

PHONY: postgres createdb dropdb migrateup migratedown test cleantest server mock migrateup1 migratedown1 dockerBuild sb dbdocs dbml2sql proto redis

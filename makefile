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
bash:
	docker exec -it post bash
access:
	docker exec -it post psql simple_bank
# sqlcdocker:
# 	docker run --rm -v "C:\Users\Gio\Documents\goworkspace\src\github.com\MikoBerries\SimpleBank":/src -w /src kjconroy/sqlc generate
# sqlcdockerver:
# 	docker run --rm -v $(PWD):/src -w //src kjconroy/sqlc version
# sqlcdockerinit:
# 	docker run --rm --volumes-from myapps -v $(PWD):/src -w /src kjconroy/sqlc init
# sqlcdockergen:
# 	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc generate
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
PHONY: postgres createdb dropdb access migrateup migratedown sqlc

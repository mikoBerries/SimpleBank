# SimpleBank
 
section database
---------------
db/migration s
 contain of file to migrade DB in golang using golang-migrate
 https://github.com/golang-migrate/migrate
 https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

$ # Go 1.16+
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$TAG

migration in sqlc
sqlc does not perform database migrations for you. However, sqlc is able to differentiate between up and down migrations. sqlc ignores down migrations when parsing SQL files.

sqlc supports parsing migrations from the following tools:
dbmate
golang-migrate
goose
sql-migrate
tern


comparing database/sql package, gorm, sqlx, and sqlc

https://blog.jetbrains.com/go/2023/04/27/comparing-db-packages/
https://gorm.io/docs/query.html

sqlc.yaml config
https://docs.sqlc.dev/en/latest/reference/config.html#gen
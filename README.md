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

Transaction 

A transaction is a single logical unit of work that accesses and possibly modifies the contents of a database. Transactions access data using read and write operations. 
In order to maintain consistency in a database, before and after the transaction, certain properties are followed. These are called ACID properties.

https://www.geeksforgeeks.org/acid-properties-in-dbms/

sql anomaly 

dirty read
lost update
non-repeatable read
phantoms
serialization anomaly
https://mkdev.me/posts/transaction-isolation-levels-with-postgresql-as-an-example


https://www.postgresql.org/docs/current/explicit-locking.html

Deadlock in database

In a database, a deadlock is a situation in which two or more transactions are waiting for one another to give up locks. For example, Transaction A might hold a lock on some rows in the Accounts table and needs to update some rows in the Orders table to finish.

Github CI

Postgres service:
https://docs.github.com/en/actions/using-containerized-services/creating-postgresql-service-containers


Gin
framework for golang 
validator used in gin
https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags
Gin ready template project
https://github.com/gin-gonic/examples
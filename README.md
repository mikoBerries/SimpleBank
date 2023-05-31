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

Dirty Read in database 

The simplest explanation of the dirty read is the state of reading uncommitted data. In this circumstance, we are not sure about the consistency of the data that is read because we don’t know the result of the open transaction(s). After reading the uncommitted data, the open transaction can be completed with rollback. On the other hand, the open transaction can complete its actions successfully. The data that is read in this ambiguous way is defined as dirty data.

To Prevent Dirty Read : Lock-Based Protocols - To attain consistency, isolation between the transactions is the most important tool. 

https://www.postgresql.org/docs/current/explicit-locking.html

Deadlock in database

In a database, a deadlock is a situation in which two or more transactions are waiting for one another to give up locks. For example, Transaction A might hold a lock on some rows in the Accounts table and needs to update some rows in the Orders table to finish.
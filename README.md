# SimpleBank
 
## section database
--------------------
1. db/migration s
 contain of file to migrade DB in golang using golang-migrate
 https://github.com/golang-migrate/migrate
 https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

$ # Go 1.16+
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$TAG

2. migration in sqlc
sqlc does not perform database migrations for you. However, sqlc is able to differentiate between up and down migrations. sqlc ignores down migrations when parsing SQL files.

sqlc supports parsing migrations from the following tools:
- dbmate
- golang-migrate
- goose
- sql-migrate
- tern


3. comparing database/sql package, gorm, sqlx, and sqlc

https://blog.jetbrains.com/go/2023/04/27/comparing-db-packages/
https://gorm.io/docs/query.html

sqlc.yaml config
https://docs.sqlc.dev/en/latest/reference/config.html#gen

4. Transaction 

A transaction is a single logical unit of work that accesses and possibly modifies the contents of a database. Transactions access data using read and write operations. 
In order to maintain consistency in a database, before and after the transaction, certain properties are followed. These are called ACID properties.

https://www.geeksforgeeks.org/acid-properties-in-dbms/

5. sql anomaly 

dirty read
lost update
non-repeatable read
phantoms
serialization anomaly
https://mkdev.me/posts/transaction-isolation-levels-with-postgresql-as-an-example


https://www.postgresql.org/docs/current/explicit-locking.html

6. Deadlock in database

In a database, a deadlock is a situation in which two or more transactions are waiting for one another to give up locks. For example, Transaction A might hold a lock on some rows in the Accounts table and needs to update some rows in the Orders table to finish.

7. Github CI

Postgres service:
https://docs.github.com/en/actions/using-containerized-services/creating-postgresql-service-containers


## section RESTFful API
-----------------------
1. Gin
framework for golang 
validator used in gin
https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags
Gin ready template project
https://github.com/gin-gonic/examples

2. Viper
Viper lib For Easy configuration File management tools
https://github.com/spf13/vigo-per

3. Mock Testing
Mock testing are using "White Box" testing techniques, used to testing every functional func in code.
https://github.com/golang/mock

4. Testing Method  TDD / BDD / ATDD
* TDD (test-driven development) approach with step:
- First creating failing unit test.
- Then start implementing unit test with MINIMAL code until all unit test condition are satisfied.
- Last testing are passed, programmer refactor the design & making an improvement without changing the behavior from second step.
in short explanation:
- TDD steps :
    - Create unit test(Code) that including all process nedeed/writen.
    - Test unitTest (Failed result test).
    - Then repeating changes/implement code until all test are passed. 
    - Done TDD cycle.
- TDD Auidence are programmer-side since it's focuses more on code implementation of a feature.
- TDD are more likely realted to "White Box" testing techniques.

* BDD (Behavioral-Driven Development) are derived/similar from TDD but using the Given-When-Then approach is used for writing test cases. Some Given-When-Then example:
    - Given the user has entered valid login credentials
    - When a user clicks on the login button
    - Then display the successful validation message
BDD in short explanation:
    - BDD Steps are similar to TDD but "unit test" test are changed to Behavioral test that writen and describe in human language (ex:English).
    - BDD Auidence are more wide to other than programmer team, because main focus on Understanding Requirements.
    - BDD are more likely realted to "Black Box" testing techniques.

 
* ATDD (Acceptance Test-Driven development) are similar with BDD focuses more on the behavior of the feature, whereas ATDD focuses on capturing the precise requirements.
ATDD in short explanation:
    - BDD are still very similar to BDD testing techniques but focusing in More detailed an precise requirements at development feature.

https://www.techtarget.com/searchsoftwarequality/definition/test-driven-development
https://www.browserstack.com/guide/tdd-vs-bdd-vs-atdd

* PASETP
paseto lib
https://paseto.io/
paseto generator online
https://token.dev/paseto/
Paserk paseto extendsion provides key-wrapping and serialization written in footer
https://github.com/paseto-standard/paserk
cross compabilty on other language
https://paseto.io/
* ETC
------ 
explanation of "var _ Interface = (*Type)(nil)"
https://github.com/uber-go/guide/issues/25
.yaml 
https://learnxinyminutes.com/docs/yaml/
Bash
https://learnxinyminutes.com/docs/bash/
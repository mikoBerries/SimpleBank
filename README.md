# SimpleBank
 
## Section Database
--------------------

1. DB Migration lib 
- contain of file to migrade DB in golang using golang-migrate
    - https://github.com/golang-migrate/migrate
    - https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
```console
$ # Go 1.16+
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$TAG
```

2. migration in sqlc
    - sqlc does not perform database migrations for you. However, sqlc is able to differentiate between up and down migrations. sqlc ignores down migrations when parsing SQL files.

    - sqlc supports parsing migrations from the following tools:
        - dbmate
        - golang-migrate
        - goose
        - sql-migrate
        - tern

3. comparing database/sql package, gorm, sqlx, and sqlc
    - article 
        - https://blog.jetbrains.com/go/2023/04/27/comparing-db-packages/
        - https://gorm.io/docs/query.html

    - sqlc.yaml config
        - https://docs.sqlc.dev/en/latest/reference/config.html#gen

4. Transaction 
    - A transaction is a single logical unit of work that accesses and possibly modifies the contents of a database. Transactions access data using read and write operations. 
    - In order to maintain consistency in a database, before and after the transaction, certain properties are followed. These are called ACID properties.

https://www.geeksforgeeks.org/acid-properties-in-dbms/

5. sql anomaly 
    - dirty read
    - lost update
    - non-repeatable read
    - phantoms
    - serialization anomaly
    https://mkdev.me/posts/transaction-isolation-levels-with-postgresql-as-an-example
    https://www.postgresql.org/docs/current/explicit-locking.html

6. Deadlock in database

    In a database, a deadlock is a situation in which two or more transactions are waiting for one another to give up locks. For example, Transaction A might hold a lock on some rows in the Accounts table and needs to update some rows in the Orders table to finish.

7. Github CI

- Postgres service:
    https://docs.github.com/en/actions/using-containerized-services/creating-postgresql-service-containers


## Section RESTFful API
-----------------------

1. Gin
- framework for golang validator used in gin
    https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags
- Gin ready template project
    https://github.com/gin-gonic/examples

2. Viper
    - Viper lib For Easy configuration File management tools (https://github.com/spf13/vigo-per)

3. Mock Testing
    - Mock testing are using "White Box" testing techniques, used to testing every functional func in code. (https://github.com/golang/mock)

4. Testing Method  TDD / BDD / ATDD
* TDD (test-driven development) approach with step:
    - First creating failing unit test.
    - Then start implementing unit test with MINIMAL code until all unit test condition are satisfied.
    - Last testing are passed, programmer refactor the design & making an improvement without changing the behavior from second step.
- TDD in short explanation:
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
* BDD in short explanation:
    - BDD Steps are similar to TDD but "unit test" test are changed to Behavioral test that writen and describe in human language (ex:English).
    - BDD Auidence are more wide to other than programmer team, because main focus on Understanding Requirements.
    - BDD are more likely realted to "Black Box" testing techniques.

 
* ATDD (Acceptance Test-Driven development) are similar with BDD focuses more on the behavior of the feature, whereas ATDD focuses on capturing the precise requirements.
ATDD in short explanation:
    - BDD are still very similar to BDD testing techniques but focusing in More detailed an precise requirements at development feature.

- https://www.techtarget.com/searchsoftwarequality/definition/test-driven-development
- https://www.browserstack.com/guide/tdd-vs-bdd-vs-atdd

5. PASETO
- paseto lib
    https://paseto.io/
- paseto generator online
    https://token.dev/paseto/
- Paserk paseto extendsion provides key-wrapping and serialization written in footer
    https://github.com/paseto-standard/paserk
- cross compabilty on other language
    https://paseto.io/

## Section Deploying aplication
-------------------------------
1. Docker images tag (alpine / bullseye / buster / windowsservercore / nanoserver):
- https://medium.com/swlh/alpine-slim-stretch-buster-jessie-bullseye-bookworm-what-are-the-differences-in-docker-62171ed4531d
2. Docker build Documentation:
- https://docs.docker.com/engine/reference/commandline/build/

3. Connecting docker container with docker network
- create a new docker network and listed images to network

4. Create docker serveral docker images with docker-compose.yaml
- With docker-compose.yaml we can automated creating serveral docker images and listed it to same network automaticly


## Section Session & GRPC
1. Session token and Access Token
    - Refresh tokens provide a way to bypass the temporary nature of access tokens. Normally, a user with an access token can only access protected resources or perform specific actions for a set period of time, which reduces the risk of the token being compromised. A refresh token allows the user to get a new access token without needing to log in again.
    - Refresh Token best practice : https://stateful.com/blog/refresh-tokens-security 

2. DbDocs.io 
    - Make sure NodeJS and NPM have been installed on your computer before the installation.
    - instal Dbdocs CLI
    ```console
    $ npm install -g dbdocs
    $ dbdocs -login
    $ dbdocs build ~/path/to/database.dbml
    $ dbdocs password --set <password> --project <project name>
    ```
    - Create your database scema in DBML in file .dbml
    - DBML (Database Markup Language) is an open-source DSL language designed to define and document database schemas and structures.
    - free website to write DBML https://dbdiagram.io/d/647f55fd722eb774947f5890
    - result https://dbdocs.io/zven_gio/simpleBank secret
    - website : https://dbdocs.io/
    - Convert a DBML file to SQL
    ```console
    $ npm install -g @dbml/cli
    $ dbml2sql schema.dbml --postgres -o schema.sql
    $ dbml2sql <path-to-dbml-file>
           [--mysql|--postgres|--mssql]
           [-o|--out-file <output-filepath>]
    ```
    - Convert a SQL file to DBML
    ```console
    $ npm install -g @dbml/cli
    $ sql2dbml dump.sql --postgres
    $ sql2dbml --mysql dump.sql -o mydatabase.dbml
    $ sql2dbml <path-to-sql-file>
           [--mysql|--postgres|--mssql]
           [-o|--out-file <output-filepath>]
    ```
    - DBML CLI make dbml-error.log file as error output even it's nothings wrong.

3. GRPC are Remote Procedure Call Framework
    - using GRPC ProtoBuff tolls can server Both GRPC and HTTP request
    - website https://grpc.io/

4. Protocol Buffer (https://protobuf.dev/programming-guides/proto3/) V3
    - protocol buffer instalation (https://grpc.io/docs/protoc-installation/) windows user using Install pre-compiled binaries (https://github.com/google/protobuf/releases)
    - for golang Install the protocol compiler plugins for Go using the following commands:
    ```console
    $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    ```
    - Protobuff extendsion vs code : https://marketplace.visualstudio.com/items?itemName=zxh404.vscode-proto3&ssr=false#overview
    - Generate protobuf from .proto to .go file
    ```console
    $ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto
    ```
    - Protobuf support scalar type and explanation (https://protobuf.dev/programming-guides/proto3/#scalar)
4. GRPC tools to testing
    - evans : https://github.com/ktr0731/evans
    ```console
    evans --host localhost -p 9090 -r repl
    ```
    - Postman also support gRPC too.
    - To use Service definition using server reflection, gRPC server mus registeret to pakcage reflection
        - reflection.Register(grpcServer)

4. GRPC VS REST api

    - https://blog.dreamfactory.com/grpc-vs-rest-how-does-grpc-compare-with-traditional-rest-apis/#:~:text=Here%20are%20the%20main%20differences,usually%20leverages%20JSON%20or%20XML.
    - https://learning.postman.com/docs/sending-requests/grpc/first-grpc-request/
5. Redis
    - Best practices https://climbtheladder.com/10-redis-key-best-practices/
## ETC
------ 
1. explanation of "var _ Interface = (*Type)(nil)"
    https://github.com/uber-go/guide/issues/25
2. .yaml 
    https://learnxinyminutes.com/docs/yaml/
3. Bash
    https://learnxinyminutes.com/docs/bash/
4. Some API management tools that support gRPC testing include Postman, Insomnia, Kreya. app, and BloomRPC
5. GRPC Status code (1 - 16) and some case code will procude:
    - https://grpc.github.io/grpc/core/md_doc_statuscodes.html
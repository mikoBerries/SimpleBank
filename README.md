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
    - Stub is replacement for some dependency in your code that will be used during test execution. It is typically built for one particular test and unlikely can be reused for another because it has hardcoded expectations and assumptions.

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

3. RPC Remote Procedure Call
    - https://medium.com/programmer-geek/mengenal-rpc-remote-procedure-call-7d8a794bbd1f
    - https://www.techtarget.com/searchapparchitecture/definition/Remote-Procedure-Call-RPC
    - Few list of RPC framework:
        1. gRPC: gRPC is a high-performance, open-source RPC framework that uses HTTP/2 as the transport layer and Protocol Buffers as the serialization format.
        2. Apache Thrift: Apache Thrift is another popular RPC framework that is available for a variety of languages.
        3. ZeroMQ: ZeroMQ is a lightweight, asynchronous RPC framework that is often used in distributed systems and microservices architectures.
        4. RabbitMQ: RabbitMQ is a popular message broker that can be used to implement RPC. 
        5. Amazon Web Services (AWS) Lambda: AWS Lambda is a serverless computing platform that can be used to implement RPC. 

4. GRPC / Google RPC Framework
    - using GRPC ProtoBuff tolls can server Both GRPC and HTTP request
    - 4 types of gRPC
        1. Unary. (1 : 1)
        2. Server Streaming. (1 : many)
        3. Client Streaming. (many : 1)
        4. Bi-directional streaming. (many : many)
    - website https://grpc.io/

5. Protocol Buffer (https://protobuf.dev/programming-guides/proto3/) V3
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
    - Protoc can produce open api documentation
    ```console
    --openapiv2_out=doc/swagger
    ```
    - Some option we can use for protoc generation write inside protobuffer services. >>>>> Do not Forget to copy .proto file to project path as import depedency
        - https://github.com/grpc-ecosystem/grpc-gateway/blob/main/examples/internal/proto/examplepb/a_bit_of_everything.proto
6. GRPC tools to testing
    - evans : https://github.com/ktr0731/evans
    ```console
    evans --host localhost -p 9090 -r repl
    ```
    - Postman also support gRPC too.
    - To use Service definition using server reflection, gRPC server mus registeret to pakcage reflection
        - reflection.Register(grpcServer)
7. gRPC gateway
    - gRPC-Gateway is a plugin of protoc. It reads a gRPC service definition and generates a reverse-proxy server which translates a RESTful JSON API into gRPC.
    - This project aims to provide that HTTP+JSON interface to your gRPC service because, you might still want to provide a traditional RESTful JSON API as well. Reasons can range from maintaining backward-compatibility, supporting languages or clients that are not well supported by gRPC
    - Source code : https://github.com/grpc-ecosystem/grpc-gateway

8. gRPC request validator
    - in gin using binding/v10 in gRPC using "google.golang.org/genproto/googleapis/rpc/errdetails"
    - when making error validtaion best practice to write field same as request param name (ex : full_name) to create consistency.

8. Authentication in GRPC
    - SSL/TLS: gRPC has SSL/TLS integration and promotes the use of SSL/TLS to authenticate the server.
    - ALTS: gRPC supports ALTS as a transport security mechanism.
    - Token-based authentication with Google : gRPC provides a generic mechanism to attach metadata based credentials to requests and responses.
        - Google credentials should only be used to connect to Google services. Sending a Google issued OAuth2 token to a non-Google service could result in this token being stolen and used to impersonate the client to Google services.
    - https://grpc.io/docs/guides/auth/
9. Logger framework for go
    - 2 good logger zap & zero log
    - https://blog.logrocket.com/5-structured-logging-packages-for-go/
    - https://github.com/rs/zerolog
    - Producing json log can help log management tolls to consume organize it.
    - https://sematext.com/blog/best-log-management-tools/

10. Swagger-ui
    - Swagger UI allows anyone — be it your development team or your end consumers — to visualize and interact with the API’s resources without having any of the implementation logic in place. It’s automatically generated from your OpenAPI (formerly known as Swagger) Specification, with the visual documentation making it easy for back end implementation and client side consumption.
    - using it in our local server with just adding swagger-ui/dist to our project and serve it as fileserver
    - https://github.com/swagger-api/swagger-ui
    - using rakyll/statik to allows you to embed a directory of static files into your Go binary (it' wil faster since it's no need to read file when called)
    - https://github.com/rakyll/statik

## Section Asynchronous processing with background workers
---------------------------------------------------------- 
1. asynq lib
    - https://github.com/hibiken/asynq still in heavy develoment release version is not done yet but high star on github
2. Redis
    - Redis is an open source (BSD licensed), in-memory data structure store used as a database, cache, message broker, and streaming engine.
    - Redis Data type https://redis.io/docs/data-types/
    - Best practices https://climbtheladder.com/10-redis-key-best-practices/
    - Always seting task with some delay so database state will be ready and redis task will satified.
    - To start redis in windows bash we must bring .conf file while starting redis-server.exe
    ```console
    $ redis-server.exe ~/path/to/redis/redis.windows.conf
    ```
3. Simple Mail Transfer Protocol (SMTP)
    - standart liblary from go :https://pkg.go.dev/net/smtp (The smtp package is frozen and is not accepting new features. Some external packages provide more functionality.)
    - lib used :https://github.com/jordan-wright/email
    - testing.Short() used flag when testing is too long to execute and we wnat to skip it (code inside .short() t.Skip())
        - adding --short flag to set it true when calling "go test"
    - AWS SES (amazon Simple Email Services)
4. Creating unit testing for gRPC services that included backgorund wokers (redis) using go-mock
    - Since there are two things need to be mocked (Database, Redis) in test code we must separate controler for each mock.
    - Code violations from database must be mapped separately to get expected error.
    - Token auth can embeded using metadata.NewIncomingContext (package google.golang.org/grpc/metadata). 

## ETC
------ 
1. explanation of "var _ Interface = (*Type)(nil)"
    - https://github.com/uber-go/guide/issues/25
2. .yaml 
    - https://learnxinyminutes.com/docs/yaml/
3. Bash
    - https://learnxinyminutes.com/docs/bash/
4. Some API management tools that support gRPC testing include Postman, Insomnia, Kreya. app, and BloomRPC
5. iana standart HTTP/2 response code
    - https://www.iana.org/assignments/http2-parameters/http2-parameters.xhtml
6. GRPC Status code (1 - 16) and some case code will procude:
    - https://grpc.github.io/grpc/core/md_doc_statuscodes.html
7. GRPC vs WebSocket
    - https://www.wallarm.com/what/grpc-vs-websocket-when-is-it-better-to-use
8. GRPC VS REST api
    - https://blog.dreamfactory.com/grpc-vs-rest-how-does-grpc-compare-with-traditional-rest-apis/#:~:text=Here%20are%20the%20main%20differences,usually%20leverages%20JSON%20or%20XML.
    - https://learning.postman.com/docs/sending-requests/grpc/first-grpc-request/
9. MQTT
    - https://aws.amazon.com/what-is/mqtt/
10. Networking Technique Multiplexing (Multiplexers and de-Multiplexers)
    - https://www.tutorialspoint.com/data_communication_computer_network/physical_layer_multiplexing.htm
11. Serialization data
    - https://hazelcast.com/glossary/serialization/
12. Anatomy of API (Application Programming Interface)
    - https://www.mertech.com/blog/the-anatomy-of-a-web-api
13. Phyton Framework Django vs Flask
    - https://www.interviewbit.com/blog/flask-vs-django/
14. Swagger hub to populate api documentation (exmaple file in doc/swagger/*.json)
    - https://swagger.io/tools/swaggerhub/
15. 5 types of design patterns:
    - Creational patterns are used to deal with creating objects.
    - Structural patterns are used to build idiomatic structures.
    - Behavioural patterns are used to manage mostly with algorithms.
    - Concurrency patterns are used to manage the timing execution and order execution of applications that have more than one flow.
        - https://blog.devgenius.io/5-useful-concurrency-patterns-in-golang-8dc90ad1ea61
    - Microservice patterns are used to build microservice applications.
16. REST API requests POST, GET, DELETE,  PUT / PATCH ?
     - https://www.abstractapi.com/guides/put-vs-patch#:~:text=PUT%20and%20PATCH%20both%20perform,entire%20body%20in%20the%20request.
     - PUT modifies a record's information and creates a new record if one is not available, and PATCH updates a resource without sending the entire body in the request.
17. Postman for testing Support http, gRPC, documenting, easy value set, etc
    - https://learning.postman.com/docs/sending-requests/variables/
18. synchronous and asynchronous
    - https://www.mendix.com/blog/asynchronous-vs-synchronous-programming/#:~:text=Asynchronous%20is%20a%20non%2Dblocking%20architecture%2C%20so%20the%20execution%20of,of%20the%20one%20before%20it.
* Default port 
    - postgre : 5432
    - redis : 6379
    - webapps : 8080
19. Docker windows/macos need to be curl in docker-machine -ip instead of localhost
    - https://forums.docker.com/t/curl-7-failed-to-connect-to-localhost-port-49160-connection-refused/7703
20. Redis docker Container connecting with redis-cli
    - https://stackoverflow.com/questions/54205691/access-redis-cli-inside-a-docker-container
21. next GRPC 
    - https://github.com/techschool/pcbook-go
    - https://www.youtube.com/playlist?list=PLy_6D98if3UJd5hxWNfAqKMr15HZqFnqf
22. kafka lib for golang
    - https://github.com/segmentio/kafka-go
    - https://hub.docker.com/r/bitnami/kafka
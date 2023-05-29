package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	sqlc "github.com/MikoBerries/SimpleBank/db/sqlc"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		log.Panic(err)
	}
	// defer db.Close()
	sqlcdb := sqlc.New(db)
	param := sqlc.CreateAccountParams{
		Owner:    "somebody",
		Balance:  10000000000000,
		Currency: "$",
	}
	result, err := sqlcdb.CreateAccount(ctx, param)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(result)
}

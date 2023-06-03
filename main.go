package main

import (
	"database/sql"
	"log"

	"github.com/MikoBerries/SimpleBank/api"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/util"

	_ "github.com/lib/pq"
)

func main() {
	//load config file using viper
	cf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	//db connection
	coon, err := sql.Open(cf.DBDriver, cf.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	//NewStore new struct of db conn and embbed querries
	store := db.NewStore(coon)
	//server config
	srv, err := api.NewServer(cf, store)
	if err != nil {
		log.Fatal(err)
	}
	// servErr := make(chan os.Signal)
	err = srv.StartServerAddress(cf.ServerAddress)

	// <-
	//FOR Grace SHUTDOWN

	// go func() {
	// 	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 		log.Printf("listen: %s\n", err)
	// 	}
	// }()

	// // Wait for interrupt signal to gracefully shutdown the server with
	// // a timeout of 5 seconds.
	// quit := make(chan os.Signal)
	// // kill (no param) default send syscall.SIGTERM
	// // kill -2 is syscall.SIGINT
	// // kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// log.Println("Shutting down server...")

	// // The context is used to inform the server it has 5 seconds to finish
	// // the request it is currently handling
	// ctx, cancel := context.WithTimeout(context.Background(), 5*coon.Stats().MaxIdleTimeClosed.Second)
	// defer cancel()

	// if err := srv.Shutdown(ctx); err != nil {
	// }

	log.Fatal("Server forced to shutdown:", err)
	log.Println("Server exiting")
}

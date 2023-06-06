package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/MikoBerries/SimpleBank/api"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/MikoBerries/SimpleBank/gapi"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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

	//uncomment this to use gin server
	//runGRPCServer(cf,store)

	//use GRPC server
	runGRPCServer(cf, store)
}

func runGRPCServer(cf util.Config, store db.Store) {
	//Make go server
	server, err := gapi.NewServer(cf, store)
	if err != nil {
		log.Fatal(err)
	}
	//make new GRPC server
	grpcServer := grpc.NewServer()

	//Register Proto Buffer to new grpc server
	pb.RegisterSimplebankServer(grpcServer, server)
	//Regiter to reflection to check whats RPC avail and how to access it
	reflection.Register(grpcServer)

	//listen tcp proctol in 0.0.0.0:9090 port
	listener, err := net.Listen("tcp", cf.GRPCServerAddress)
	if err != nil {
		log.Fatal("error when setting gRPC server : ", err)
	}
	log.Printf("Starting gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("error when start gRPC server : ", err)
	}
}

// runGinServer func to start server as HTTP using Gin
func runGinServer(cf util.Config, store db.Store) {
	srv, err := api.NewServer(cf, store)
	if err != nil {
		log.Fatal(err)
	}
	// servErr := make(chan os.Signal)
	err = srv.StartServerAddress(cf.HttpServerAddress)

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

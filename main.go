package main

import (
	"context"
	"database/sql"

	"net"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/MikoBerries/SimpleBank/api"
	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	_ "github.com/MikoBerries/SimpleBank/doc/statik"
	"github.com/MikoBerries/SimpleBank/gapi"
	"github.com/MikoBerries/SimpleBank/pb"
	"github.com/MikoBerries/SimpleBank/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	//load config file using viper
	cf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
	if cf.Enviroment != "Production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	//db connection
	coon, err := sql.Open(cf.DBDriver, cf.DBSource)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	//migrate database schema
	migrateDatabase(cf.DBMigratePath, cf.DBSource)

	//NewStore new struct of db conn and embbed querries
	store := db.NewStore(coon)

	//uncomment this to use gin server
	//runGRPCServer(cf,store)

	//run HTPP proxy server
	go runHTTPServer(cf, store)
	//use GRPC server
	runGRPCServer(cf, store)
}

// runHTTPServer are gPRC Proxy server to serve Http json and forwading to gRPC server
func runHTTPServer(cf util.Config, store db.Store) {
	server, err := gapi.NewServer(cf, store)
	if err != nil {
		log.Fatal().Msgf("err: %s", err)
	}
	ctx := context.Background()
	//embed context with cancel
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//multiplexer server

	gatewayOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{ //marshal func option
			UseProtoNames: true, // use callback name from .proto file
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true, //discard eveeything unknow field
		},
	})

	grpcMux := runtime.NewServeMux(gatewayOption)

	// err = pb.RegisterSimplebankHandlerFromEndpoint(ctx, grpcMux, server, nil)
	//register handler path from gAPI to grpc mux
	err = pb.RegisterSimplebankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msgf(err.Error())
	}

	mux := http.NewServeMux()
	//handle all path
	mux.Handle("/", grpcMux)

	//make swagger ui statik file
	statikFile, err := fs.New()
	if err != nil {
		log.Fatal().Msgf("err:", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFile))
	mux.Handle("/swagger/", swaggerHandler)

	//lsitener to listen Tcp in 8080 port
	listener, err := net.Listen("tcp", cf.HttpServerAddress)
	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
	log.Printf("Starting HTTP Proxy server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
}

func runGRPCServer(cf util.Config, store db.Store) {
	//Make go server
	server, err := gapi.NewServer(cf, store)
	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
	//logging for unary traffic
	opt := grpc.UnaryInterceptor(gapi.GRPCLogger)
	//make new GRPC server
	grpcServer := grpc.NewServer(opt)

	//Register Proto Buffer to new grpc server
	pb.RegisterSimplebankServer(grpcServer, server)
	//Regiter to reflection to check whats RPC avail and how to access it
	reflection.Register(grpcServer)

	//listen tcp proctol in 0.0.0.0:9090 port
	listener, err := net.Listen("tcp", cf.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msgf("error when setting gRPC server : %s", err.Error())
	}
	log.Printf("Starting gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msgf("error when start gRPC server : %s", err.Error())
	}
}

// runGinServer func to start server as HTTP using Gin
func runGinServer(cf util.Config, store db.Store) {
	srv, err := api.NewServer(cf, store)
	if err != nil {
		log.Fatal().Msgf("err:", err)
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

	log.Fatal().Msgf("Server forced to shutdown:%s", err.Error())
	log.Print("Server exiting")
}

// MigrateDatabase to execute mirgate database before server starting
func migrateDatabase(migratePath string, DBSource string) {
	log.Print("Start migrate database migrate")

	migration, err := migrate.New(migratePath, DBSource)
	if err != nil {
		log.Fatal().Msgf("error when setting migrate :%s", err.Error())
	}
	//even with migrate run well will returning err no change
	if err = migration.Up(); err != migrate.ErrNoChange {
		log.Fatal().Msgf("error when migrate database :%s", err.Error())
	}

	log.Print("Done migrate database migrate")
}

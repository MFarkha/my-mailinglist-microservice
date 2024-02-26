package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/MFarkha/my-mailinglist-microservice/grpcapi"
	"github.com/MFarkha/my-mailinglist-microservice/jsonapi"
	"github.com/MFarkha/my-mailinglist-microservice/mdb"
	"go.wit.com/dev/alexflint/arg"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"` // address:port to listen JSON server
	BindgRPC string `arg:"env:MAILINGLIST_BIND_GRPC"` // address:port to listen for gRPC server
}

func main() {
	arg.MustParse(&args)
	if args.DbPath == "" {
		args.DbPath = "_data/list.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":3000"
	}
	if args.BindgRPC == "" {
		args.BindgRPC = ":3001"
	}
	log.Printf("Using the database: '%v'\n", args.DbPath)
	db, err := sql.Open("sqlite3", args.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mdb.TryCreate(db)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("starting JSON API server...\n")
		jsonapi.Serve(db, args.BindJson)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("starting gRPC API server...\n")
		grpcapi.Serve(db, args.BindgRPC)
	}()
	wg.Wait()
}

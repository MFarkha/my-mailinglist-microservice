package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/MFarkha/my-mailinglist-microservice/jsonapi"
	"github.com/MFarkha/my-mailinglist-microservice/mdb"
	"go.wit.com/dev/alexflint/arg"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"` // address:port to listen
}

func main() {
	arg.MustParse(&args)
	if args.DbPath == "" {
		args.DbPath = "_data/list.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":3000"
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
		log.Printf("The app is listening the port%s\n", args.BindJson)
		jsonapi.Serve(db, args.BindJson)
	}()
	wg.Wait()
}

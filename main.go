package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/enrichman/api-fosdem/indexer"
	"github.com/enrichman/api-fosdem/store"
)

func main() {
	port := os.Getenv("PORT")

	mongoStore, err := store.NewMongoStore("", "")
	if err != nil {
		panic(err)
	}
	remoteIndexer := indexer.NewRemoteIndexer(mongoStore)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/reindex", indexer.MakeIndexerHandler(remoteIndexer))
	http.Handle("/", mux)

	fmt.Println("listening...", port)

	srv := http.Server{Addr: ":" + port}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}

	fmt.Println("closed.")
}

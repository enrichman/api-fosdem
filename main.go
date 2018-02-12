package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/enrichman/api-fosdem/indexer"
)

func main() {
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.Handle("/api/v1/reindex", indexer.MakeIndexerHandler(&indexer.RemoteIndexer{}))
	http.Handle("/", mux)

	fmt.Println("listening...", port)

	srv := http.Server{Addr: ":" + port}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}

	fmt.Println("closed.")
}

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/enrichman/api-fosdem/indexer"
	"github.com/enrichman/api-fosdem/speakers"
	"github.com/enrichman/api-fosdem/store"
)

func main() {
	port := os.Getenv("PORT")
	token := os.Getenv("TOKEN")
	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")

	mongoStore, err := store.NewMongoStore(mongoURI, mongoDB)
	if err != nil {
		panic(err)
	}
	remoteIndexer := &indexer.RemoteIndexer{
		Token: token, SpeakerSaver: mongoStore,
	}
	remoteFinder := &speakers.RemoteGetter{SpeakerFinder: mongoStore}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/reindex", indexer.MakeIndexerHandler(remoteIndexer))
	mux.Handle("/api/v1/", speakers.MakeSpeakersHandler(remoteFinder))
	http.Handle("/", mux)

	fmt.Println("listening...", port)

	srv := http.Server{Addr: ":" + port}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}

	fmt.Println("closed.")
}

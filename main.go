package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/enrichman/api-fosdem/indexer"
	"github.com/enrichman/api-fosdem/pentabarf"
	"github.com/enrichman/api-fosdem/speakers"
	"github.com/enrichman/api-fosdem/store"
	"github.com/enrichman/api-fosdem/web"
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
	remoteIndexer := indexer.NewRemoteIndexer(
		token,
		&pentabarf.CachedScheduleService{},
		mongoStore,
		web.NewSpeakerService(),
	)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/reindex", indexer.MakeReindexerHandler(remoteIndexer))
	mux.Handle("/api/v1/", speakers.MakeSpeakersHandler(speakers.NewService(mongoStore)))
	http.Handle("/", mux)

	fmt.Println("listening...", port)

	srv := http.Server{Addr: ":" + port}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}

	fmt.Println("closed.")
}

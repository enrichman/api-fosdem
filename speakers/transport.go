package speakers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeSpeakersHandler(s SpeakerFinder) http.Handler {
	r := mux.NewRouter()

	speakerGetterHandler := kithttp.NewServer(
		makeSpeakerGetterEndpoint(s),
		decodeSpeakerGetter,
		encodeSpeakerGetter,
	)

	speakerFinderHandler := kithttp.NewServer(
		makeSpeakerFinderEndpoint(s),
		decodeSpeakerFinder,
		encodeSpeakerFinder,
	)

	r.Handle("/api/v1/speakers", speakerFinderHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/speakers/{id}", speakerGetterHandler).Methods(http.MethodGet)

	return r
}

func decodeSpeakerGetter(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, errors.New("wrong ID")
	}
	return id, nil
}

func encodeSpeakerGetter(_ context.Context, w http.ResponseWriter, res interface{}) error {
	return json.NewEncoder(w).Encode(res)
}

func decodeSpeakerFinder(_ context.Context, r *http.Request) (request interface{}, err error) {
	return r.FormValue("name"), nil
}

func encodeSpeakerFinder(_ context.Context, w http.ResponseWriter, res interface{}) error {
	return json.NewEncoder(w).Encode(res)
}

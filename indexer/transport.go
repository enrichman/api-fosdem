package indexer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeReindexerHandler setup the handlers on the /api/v1/reindex route
func MakeReindexerHandler(i indexer) http.Handler {
	r := mux.NewRouter()

	reindexHandler := kithttp.NewServer(
		makeReindexEndpoint(i),
		decodeReindex,
		encodeReindex,
	)

	r.Handle("/api/v1/reindex", reindexHandler).Methods(http.MethodGet)

	return r
}

func decodeReindex(_ context.Context, r *http.Request) (request interface{}, err error) {
	token := r.FormValue("token")
	if token == "" {
		return nil, errors.New("missing token")
	}
	return reindexRequest{token}, nil
}

func encodeReindex(_ context.Context, w http.ResponseWriter, res interface{}) error {
	resp := res.(reindexResponse)
	return json.NewEncoder(w).Encode(resp)
}

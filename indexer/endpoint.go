package indexer

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type reindexRequest struct {
	token string
}

type reindexResponse struct {
	Err error `json:"error"`
}

func makeReindexEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return reindexResponse{Err: errors.New("ciao")}, nil
	}
}

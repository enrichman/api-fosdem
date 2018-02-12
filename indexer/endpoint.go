package indexer

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type reindexRequest struct {
	token string
}

type reindexResponse struct {
	Err error `json:"error"`
}

type Indexer interface {
	Index() error
}

func makeReindexEndpoint(indexer Indexer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return reindexResponse{indexer.Index()}, nil
	}
}

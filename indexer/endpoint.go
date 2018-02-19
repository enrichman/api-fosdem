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

type indexer interface {
	GetToken() string
	Index() error
}

func makeReindexEndpoint(indexer indexer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(reindexRequest)
		if req.token != indexer.GetToken() {
			return nil, errors.New("invalid token")
		}
		return reindexResponse{indexer.Index()}, nil
	}
}

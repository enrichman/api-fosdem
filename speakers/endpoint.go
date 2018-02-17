package speakers

import (
	"context"

	"github.com/enrichman/api-fosdem/indexer"
	"github.com/go-kit/kit/endpoint"
)

type SpeakerFinder interface {
	FindByID(int) (*indexer.Speaker, error)
	Find(limit, offset int, name string, years []int) ([]indexer.Speaker, int, error)
}

func makeSpeakerGetterEndpoint(finder SpeakerFinder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return finder.FindByID(request.(int))
	}
}

type findRequest struct {
	limit  int
	offset int
	slug   string
	years  []int
}

type findResponse struct {
	Count int               `json:"count"`
	Data  []indexer.Speaker `json:"data"`
}

func makeSpeakerFinderEndpoint(finder SpeakerFinder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(findRequest)
		speakers, count, err := finder.Find(req.limit, req.offset, req.slug, req.years)
		if err != nil {
			return nil, err
		}
		return findResponse{
			Count: count,
			Data:  speakers,
		}, nil
	}
}

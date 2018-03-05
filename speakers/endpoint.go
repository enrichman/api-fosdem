package speakers

import (
	"context"

	"github.com/enrichman/api-fosdem/store"
	"github.com/go-kit/kit/endpoint"
)

type speakerFinder interface {
	FindByID(int) (*store.Speaker, error)
	Find(limit, offset int, name string, years []int) ([]store.Speaker, int, error)
}

func makeSpeakerGetterEndpoint(finder speakerFinder) endpoint.Endpoint {
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
	Count int             `json:"count"`
	Data  []store.Speaker `json:"data"`
}

func makeSpeakerFinderEndpoint(finder speakerFinder) endpoint.Endpoint {
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

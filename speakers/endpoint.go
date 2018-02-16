package speakers

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeSpeakerGetterEndpoint(finder SpeakerFinder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return finder.FindByID(request.(int))
	}
}

func makeSpeakerFinderEndpoint(finder SpeakerFinder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return finder.Find(request.(string)), nil
	}
}

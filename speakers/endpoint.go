package speakers

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type speakerService interface {
	FindByID(int) (*Speaker, error)
	Find(limit, offset int, name string, years []int) ([]Speaker, int, error)
}

func makeSpeakerGetterEndpoint(finder speakerService) endpoint.Endpoint {
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
	Count int       `json:"count"`
	Data  []Speaker `json:"data"`
}

// Speaker maps the speaker
type Speaker struct {
	ID           int    `json:"id,omitempty"`
	Slug         string `json:"slug,omitempty"`
	Name         string `json:"name,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
	ProfilePage  string `json:"profile_page,omitempty"`
	Bio          string `json:"bio,omitempty"`
	Year         int    `json:"year,omitempty"`
	Links        []Link `json:"links,omitempty"`
}

// Link is a detail link owned by a Speaker
type Link struct {
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}

func makeSpeakerFinderEndpoint(finder speakerService) endpoint.Endpoint {
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

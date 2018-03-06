package speakers

import "github.com/enrichman/api-fosdem/store"

type speakerFinder interface {
	FindByID(int) (*store.Speaker, error)
	Find(limit, offset int, name string, years []int) ([]store.Speaker, int, error)
}

type Service struct {
	speakerFinder speakerFinder
}

func NewService(speakerFinder speakerFinder) *Service {
	return &Service{speakerFinder}
}

func (s *Service) FindByID(id int) (*Speaker, error) {
	storeSpeaker, err := s.speakerFinder.FindByID(id)
	if err != nil {
		return nil, err
	}
	speaker := convertSpeaker(*storeSpeaker)
	return &speaker, nil
}

func (s *Service) Find(limit, offset int, name string, years []int) ([]Speaker, int, error) {
	speakersFound, count, err := s.speakerFinder.Find(limit, offset, name, years)
	if err != nil {
		return nil, -1, err
	}

	speakers := make([]Speaker, 0)
	for _, s := range speakersFound {
		speakers = append(speakers, convertSpeaker(s))
	}

	return speakers, count, nil
}

func convertSpeaker(s store.Speaker) Speaker {
	speaker := Speaker{
		ID:           s.ID,
		Slug:         s.Slug,
		Name:         s.Name,
		ProfileImage: s.ProfileImage,
		ProfilePage:  s.ProfilePage,
		Bio:          s.Bio,
		Year:         s.Year,
		Links:        make([]Link, 0),
	}
	for _, l := range s.Links {
		speaker.Links = append(speaker.Links, Link{URL: l.URL, Title: l.Title})
	}
	return speaker
}

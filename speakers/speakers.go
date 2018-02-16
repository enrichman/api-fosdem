package speakers

import (
	"github.com/enrichman/api-fosdem/indexer"
)

const (
	baseUrl = "https://fosdem.org"
)

type SpeakerFinder interface {
	FindByID(int) (*indexer.Speaker, error)
	Find(string) []indexer.Speaker
}

type RemoteGetter struct {
	SpeakerFinder SpeakerFinder
}

func (fi *RemoteGetter) FindByID(ID int) (*indexer.Speaker, error) {
	return fi.SpeakerFinder.FindByID(ID)
}

func (fi *RemoteGetter) Find(query string) []indexer.Speaker {
	return fi.SpeakerFinder.Find(query)
}

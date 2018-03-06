package indexer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/enrichman/api-fosdem/pentabarf"

	"github.com/enrichman/api-fosdem/store"
	"github.com/enrichman/api-fosdem/web"
)

type speakerSaver interface {
	Save(s store.Speaker) error
}

type scheduleGetter interface {
	GetSchedule(year int) (*pentabarf.Schedule, error)
}

type speakerGetter interface {
	GetSpeakers() <-chan web.Result
	GetSpeakersByYear(int) <-chan web.Result
}

// RemoteIndexer is an indexer that fetch the FOSDEM XML remotely
type RemoteIndexer struct {
	Token          string
	scheduleGetter scheduleGetter
	speakerSaver   speakerSaver
	speakerGetter  speakerGetter
}

// NewRemoteIndexer returns a remoteIndexer
func NewRemoteIndexer(
	token string,
	scheduleGetter scheduleGetter,
	speakerSaver speakerSaver,
	speakerGetter speakerGetter,
) *RemoteIndexer {
	return &RemoteIndexer{
		Token:          token,
		scheduleGetter: scheduleGetter,
		speakerSaver:   speakerSaver,
		speakerGetter:  speakerGetter,
	}
}

// GetToken returns the token used to check if the request is valid
func (fi *RemoteIndexer) GetToken() string {
	return fi.Token
}

// Index starts the indexing
func (fi *RemoteIndexer) Index() error {
	start := time.Now()
	fmt.Println(start, "start indexing")

	for year := 2013; year <= 2018; year++ {
		err := fi.IndexYear(year)
		if err != nil {
			fmt.Println("error indexing year " + strconv.Itoa(year))
		}
	}

	fmt.Println(time.Since(start), "finished indexing")
	return nil
}

// IndexYear index the provided year
func (fi *RemoteIndexer) IndexYear(year int) error {
	schedule, err := fi.scheduleGetter.GetSchedule(year)
	if err != nil {
		return err
	}

	count := 0
	for r := range fi.speakerGetter.GetSpeakersByYear(year) {
		if r.Error != nil {
			fmt.Println("error getting speaker: " + r.Error.Error())
			continue
		}

		p, found := schedule.GetPersonByName(r.Speaker.Name)
		if !found {
			fmt.Println("person with name not found in schedule: " + r.Speaker.Name)
			continue
		}

		s := store.Speaker{
			ID:           p.ID,
			Slug:         r.Speaker.Slug,
			Name:         r.Speaker.Name,
			ProfileImage: r.Speaker.ProfileImage,
			ProfilePage:  r.Speaker.ProfilePage,
			Bio:          r.Speaker.Bio,
			Year:         r.Speaker.Year,
		}

		s.Links = make([]store.Link, 0)
		for _, l := range r.Speaker.Links {
			s.Links = append(s.Links, store.Link{Title: l.Title, URL: l.URL})
		}

		err = fi.speakerSaver.Save(s)
		if err != nil {
			fmt.Println("error saving speaker: " + err.Error())
			continue
		}

		count++
		fmt.Printf("%d) year [%d] speaker [%s] saved\n", count, year, s.Name)
	}

	return nil
}

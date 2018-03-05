package indexer

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/enrichman/api-fosdem/pentabarf"

	"github.com/enrichman/api-fosdem/store"
	"github.com/enrichman/api-fosdem/web"
)

const (
	baseURL = "https://fosdem.org"
)

type speakerSaver interface {
	Save(s store.Speaker) error
}

type speakerGetter interface {
	GetSpeakers() <-chan web.Result
	GetSpeakersByYear(int) <-chan web.Result
}

// RemoteIndexer is an indexer that fetch the FOSDEM XML remotely
type RemoteIndexer struct {
	Token         string
	speakerSaver  speakerSaver
	speakerGetter speakerGetter
}

// NewRemoteIndexer returns a remoteIndexer
func NewRemoteIndexer(token string, speakerSaver speakerSaver) *RemoteIndexer {
	return &RemoteIndexer{
		Token:         token,
		speakerSaver:  speakerSaver,
		speakerGetter: web.NewSpeakerService(),
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
		scheduleRes, err := http.Get(baseURL + "/" + strconv.Itoa(year) + "/schedule/xml")
		if err != nil {
			return err
		}

		schedule, err := pentabarf.Parse(scheduleRes.Body)
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
	}

	fmt.Println(time.Since(start), "finished indexing")
	return nil
}

// IndexYear index the provided year
func (fi *RemoteIndexer) IndexYear(year string) error {
	fmt.Println("start indexing of year " + year)
	start := time.Now()

	scheduleRes, err := http.Get(baseURL + "/" + year + "/schedule/xml")
	if err != nil {
		return err
	}

	speakers, err := ParseScheduleXML(scheduleRes.Body)
	if err != nil {
		return err
	}

	speakersRes, err := http.Get(baseURL + "/" + year + "/schedule/speakers/")
	if err != nil {
		return err
	}
	speakersLink := ParseSpeakersPage(speakersRes.Body)

	speakers = FillSpeakersInfo(speakers, speakersLink)
	for _, s := range speakers {
		resp, err := http.Get(baseURL + s.ProfilePage)
		if err == nil {
			err := ParseSpeakerPage(&s, resp.Body)
			if err == nil {
				s.Slug = getSlugByLink(s.ProfilePage)
				_, _ = strconv.Atoi(year)
				/*err := fi.speakerSaver.Save(s, y)
				if err != nil {
					fmt.Println("error upserting:", s.Name)
				}*/
			}
		}
	}
	fmt.Printf("end indexing of year %s in %+v. Indexed %d speakers.\n", year, time.Since(start), len(speakers))

	return nil
}

func getSlugByLink(detailLink string) string {
	if detailLink == "" {
		return ""
	}
	arr := strings.Split(detailLink, "/")
	cleanedArr := make([]string, 0)
	for _, s := range arr {
		if len(s) > 0 {
			cleanedArr = append(cleanedArr, s)
		}
	}
	if len(cleanedArr) == 0 {
		return ""
	}
	return cleanedArr[len(cleanedArr)-1]
}

// LocalIndexer is an indexer that fetch the FOSDEM XML locally
type LocalIndexer struct{}

// Index starts the indexing
func (fi *LocalIndexer) Index() error {
	speakersFile, err := os.Open("schedule.xml")
	if err != nil {
		return err
	}

	speakers, err := ParseScheduleXML(speakersFile)
	if err != nil {
		return err
	}

	speakersPage, err := os.Open("fosdem-speakers.htm")
	if err != nil {
		return err
	}
	speakersLink := ParseSpeakersPage(speakersPage)

	speakers = FillSpeakersInfo(speakers, speakersLink)
	speakers = speakers[0:5]

	for _, s := range speakers {
		fmt.Println("Getting", s.Name)

		resp, err := http.Get(baseURL + s.ProfilePage)
		if err == nil {
			err := ParseSpeakerPage(&s, resp.Body)
			if err == nil {
				fmt.Println(s)
			}
		}
	}
	return nil
}

package indexer

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL = "https://fosdem.org"
)

type speakerSaver interface {
	Save(s Speaker, year int) error
}

// RemoteIndexer is an indexer that fetch the FOSDEM XML remotely
type RemoteIndexer struct {
	Token        string
	speakerSaver speakerSaver
}

// NewRemoteIndexer returns a remoteIndexer
func NewRemoteIndexer(token string, speakerSaver speakerSaver) *RemoteIndexer {
	return &RemoteIndexer{
		Token:        token,
		speakerSaver: speakerSaver,
	}
}

// GetToken returns the token used to check if the request is valid
func (fi *RemoteIndexer) GetToken() string {
	return fi.Token
}

// Index starts the indexing
func (fi *RemoteIndexer) Index() error {
	years := []string{"2013", "2014", "2015", "2016", "2017", "2018"}
	for _, y := range years {
		err := fi.IndexYear(y)
		if err != nil {
			return err
		}
	}
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
	speakers := ParseScheduleXML(scheduleRes.Body)

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
				y, _ := strconv.Atoi(year)
				err := fi.speakerSaver.Save(s, y)
				if err != nil {
					fmt.Println("error upserting:", s.Name)
				}
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

type localIndexer struct{}

func (fi *localIndexer) Index() error {
	speakersFile, err := os.Open("schedule.xml")
	if err != nil {
		return err
	}
	speakers := ParseScheduleXML(speakersFile)

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

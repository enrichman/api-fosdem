package indexer

import (
	"fmt"
	"net/http"
	"os"
)

const (
	baseUrl = "https://fosdem.org"
)

type SpeakerSaver interface {
	Save(Speaker) error
}

type RemoteIndexer struct {
	speakerSaver SpeakerSaver
}

func NewRemoteIndexer(speakerSaver SpeakerSaver) *RemoteIndexer {
	return &RemoteIndexer{speakerSaver}
}

func (fi *RemoteIndexer) Index() error {
	scheduleRes, err := http.Get(baseUrl + "/2018/schedule/xml")
	if err != nil {
		return err
	}
	speakers := ParseScheduleXML(scheduleRes.Body)

	speakersRes, err := http.Get(baseUrl + "/2018/schedule/speakers/")
	if err != nil {
		return err
	}
	speakersLink := ParseSpeakersPage(speakersRes.Body)

	speakers = FillSpeakersInfo(speakers, speakersLink)
	for _, s := range speakers {
		resp, err := http.Get(baseUrl + s.ProfilePage)
		if err == nil {
			err := ParseSpeakerPage(&s, resp.Body)
			if err == nil {
				err := fi.speakerSaver.Save(s)
				if err == nil {
					fmt.Println("upserted", s.Name)
				}
			}
		}
	}
	return nil
}

type LocalIndexer struct{}

func (fi *LocalIndexer) Index() error {
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

		resp, err := http.Get(baseUrl + s.ProfilePage)
		if err == nil {
			err := ParseSpeakerPage(&s, resp.Body)
			if err == nil {
				fmt.Println(s)
			}
		}
	}

	return nil
}

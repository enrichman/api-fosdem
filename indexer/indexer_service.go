package indexer

import (
	"fmt"
	"net/http"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	baseUrl = "https://fosdem.org"
)

type RemoteIndexer struct{}

func (fi *RemoteIndexer) Index() error {
	start := time.Now()

	session, err := mgo.Dial("")
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB("api-fosdem").C("speakers")

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
				_, err := c.Upsert(bson.M{"_id": s.ID}, s)
				if err == nil {
					fmt.Println("upserted", s.Name)
				}
			}
		}
	}

	fmt.Println("done", time.Since(start))

	return nil
}

type LocalIndexer struct{}

func (fi *LocalIndexer) Index() error {
	session, err := mgo.Dial("")
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB("api-fosdem").C("speakers")

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
				chInfo, err := c.Upsert(bson.M{"_id": s.ID}, s)
				fmt.Println("upserted", chInfo, err)
			}
		}
	}

	return nil
}

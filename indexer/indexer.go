package indexer

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Schedule struct {
	Days []Day `xml:"day"`
}

type Day struct {
	Index int    `xml:"index,attr"`
	Date  string `xml:"date,attr"`
	Rooms []Room `xml:"room"`
}

type Room struct {
	Name   string  `xml:"name,attr"`
	Events []Event `xml:"event"`
}

type Event struct {
	ID       int       `xml:"id,attr"`
	Speakers []Speaker `xml:"persons>person"`
}

type Speaker struct {
	ID           string `json:"id" xml:"id,attr"`
	Name         string `json:"name" xml:",chardata"`
	ProfileImage string `json:"profile_image" xml:"-"`
	ProfilePage  string `json:"profile_page" xml:"-"`
	Bio          string `json:"bio" xml:"-"`
	Links        []Link `json:"links" xml:"-"`
}

type Link struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func ParseScheduleXML(xmlReader io.Reader) []Speaker {
	var schedule Schedule
	xml.NewDecoder(xmlReader).Decode(&schedule)

	speakers := make([]Speaker, 0)
	for _, d := range schedule.Days {
		for _, r := range d.Rooms {
			for _, ev := range r.Events {
				for _, s := range ev.Speakers {
					speakers = append(speakers, s)
				}
			}
		}
	}

	return speakers
}

//ParseSpeakersPage returns a map SpeakerName to DetailPageLink of the speakers
func ParseSpeakersPage(htmlPage io.Reader) map[string]string {
	root, err := html.Parse(htmlPage)
	if err != nil {
		panic(err)
	}

	speakersLinks := scrape.FindAll(root, func(n *html.Node) bool {
		if n.DataAtom != atom.A {
			return false
		}
		return strings.Contains(scrape.Attr(n, "href"), "/schedule/speaker/")
	})

	linkMap := make(map[string]string)
	noTrimJoiner := func(s []string) string { return strings.Join(s, "") }

	for _, link := range speakersLinks {
		name := scrape.TextJoin(link, noTrimJoiner)
		linkMap[name] = scrape.Attr(link, "href")
	}

	return linkMap
}

func FillSpeakersInfo(speakers []Speaker, detailLinkMap map[string]string) []Speaker {
	fullSpeakers := make([]Speaker, 0)

	for _, s := range speakers {
		detailLink, found := detailLinkMap[s.Name]
		if !found {
			fmt.Println("detail link not found for:", s.Name)
		} else {
			s.ProfilePage = detailLink
			fullSpeakers = append(fullSpeakers, s)
		}
	}

	return fullSpeakers
}

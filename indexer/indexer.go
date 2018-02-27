package indexer

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/enrichman/api-fosdem/pentabarf"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Speaker maps the speaker
type Speaker struct {
	ID           int    `json:"id" bson:"_id"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image,omitempty"`
	ProfilePage  string `json:"profile_page"`
	Bio          string `json:"bio,omitempty"`
	Links        []Link `json:"links,omitempty"`
	Years        []int  `json:"years,omitempty"`
}

// Link is a detail link owned by a Speaker
type Link struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// ParseScheduleXML parse the XML returning the list of speakers
func ParseScheduleXML(xmlReader io.Reader) ([]Speaker, error) {
	schedule, err := pentabarf.Parse(xmlReader)
	if err != nil {
		return nil, err
	}

	speakers := make([]Speaker, 0)
	for _, ev := range schedule.GetAllEvents() {
		for _, p := range ev.Persons {
			speakers = append(speakers, Speaker{
				ID: p.ID, Name: p.Name,
			})
		}
	}

	return speakers, nil
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

// FillSpeakersInfo fills the speakers information retrieving the info from the urls passed in the map
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

// ParseSpeakerPage parse the Speaker HTML page and fills the Speaker info
func ParseSpeakerPage(speaker *Speaker, htmlPage io.Reader) error {
	root, err := html.Parse(htmlPage)
	if err != nil {
		return err
	}

	mainDiv, found := scrape.Find(root, func(n *html.Node) bool {
		return scrape.Attr(n, "id") == "main"
	})
	if !found {
		return errors.New("main div not found")
	}

	// <p> data and print them
	bioDatas := scrape.FindAll(mainDiv, pMatcher)
	for _, bioData := range bioDatas {
		speaker.Bio += scrape.Text(bioData) + "\n\n"
	}
	speaker.Bio = strings.TrimSpace(speaker.Bio)

	if imgNode, found := scrape.Find(mainDiv, imgMatcher); found {
		speaker.ProfileImage = scrape.Attr(imgNode, "src")
	}

	if h1Node, found := scrape.Find(mainDiv, h1Matcher); found {
		speaker.Name = scrape.Text(h1Node)
	}

	speaker.Links = make([]Link, 0)
	for _, h3Node := range scrape.FindAll(mainDiv, h3Matcher) {
		if scrape.Text(h3Node) == "Links" {
			if ulNode, foundUl := scrape.FindNextSibling(h3Node, ulMatcher); foundUl {
				for _, aNode := range scrape.FindAll(ulNode, aMatcher) {
					speaker.Links = append(speaker.Links, Link{
						Title: scrape.Text(aNode),
						URL:   scrape.Attr(aNode, "href"),
					})
				}
			}
		}
	}

	return nil
}

func h1Matcher(n *html.Node) bool  { return n.DataAtom == atom.H1 }
func h3Matcher(n *html.Node) bool  { return n.DataAtom == atom.H3 }
func aMatcher(n *html.Node) bool   { return n.DataAtom == atom.A }
func pMatcher(n *html.Node) bool   { return n.DataAtom == atom.P }
func imgMatcher(n *html.Node) bool { return n.DataAtom == atom.Img }
func ulMatcher(n *html.Node) bool  { return n.DataAtom == atom.Ul }

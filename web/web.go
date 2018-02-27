package web

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	baseURL      = "https://fosdem.org"
	pathSpeakers = "/schedule/speakers/"
)

type Speaker struct {
	Name        string
	ProfilePage string
}

func GetSpeakers() {

}

func getSpeakers(year int) ([]Speaker, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d%s", baseURL, year, pathSpeakers))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseSpeakers(resp.Body)
}

//ParseSpeakersPage returns a map SpeakerName to DetailPageLink of the speakers
func parseSpeakers(htmlPage io.Reader) ([]Speaker, error) {
	root, err := html.Parse(htmlPage)
	if err != nil {
		return nil, err
	}

	speakersLinks := scrape.FindAll(root, func(n *html.Node) bool {
		if n.DataAtom != atom.A {
			return false
		}
		return strings.Contains(scrape.Attr(n, "href"), "/schedule/speaker/")
	})

	speakers := make([]Speaker, 0)
	for _, link := range speakersLinks {
		speakers = append(speakers, Speaker{
			Name:        scrape.TextJoin(link, noTrimJoiner),
			ProfilePage: scrape.Attr(link, "href"),
		})
	}
	return speakers, nil
}

func getSpeaker(profilePage string) (*Speaker, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, profilePage))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseSpeaker(resp.Body)
}

func parseSpeaker(htmlPage io.Reader) (*Speaker, error) {
	root, err := html.Parse(htmlPage)
	if err != nil {
		return nil, err
	}

	mainDiv, found := scrape.Find(root, func(n *html.Node) bool {
		return scrape.Attr(n, "id") == "main"
	})
	if !found {
		return nil, errors.New("main div not found")
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

	return nil, nil
}

func noTrimJoiner(s []string) string { return strings.Join(s, "") }

func h1Matcher(n *html.Node) bool  { return n.DataAtom == atom.H1 }
func h3Matcher(n *html.Node) bool  { return n.DataAtom == atom.H3 }
func aMatcher(n *html.Node) bool   { return n.DataAtom == atom.A }
func pMatcher(n *html.Node) bool   { return n.DataAtom == atom.P }
func imgMatcher(n *html.Node) bool { return n.DataAtom == atom.Img }
func ulMatcher(n *html.Node) bool  { return n.DataAtom == atom.Ul }

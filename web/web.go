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

var years = []int{2018}

type Speaker struct {
	Name         string
	Bio          string
	ProfilePage  string
	ProfileImage string
	Year         int
	Links        []Link
}

type Link struct {
	Title string
	URL   string
}

type Result struct {
	Speaker Speaker
	Error   error
}

type speakerGetter interface {
	GetSpeakersByYear(year int) (io.Reader, error)
	GetSpeaker(profilePage string) (io.Reader, error)
}

type remoteGetter struct{}

func (g *remoteGetter) GetSpeakersByYear(year int) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d%s", baseURL, year, pathSpeakers))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Body, nil
}

func (g *remoteGetter) GetSpeaker(profilePage string) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, profilePage))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Body, nil
}

type WebSpeakerService struct {
	g speakerGetter
}

func (w *WebSpeakerService) GetSpeakers() <-chan Result {
	c := make(chan Result)

	go func() {
		for _, y := range years {
			reader, err := w.g.GetSpeakersByYear(y)
			if err != nil {
				c <- Result{Error: err}
				continue
			}

			speakers, err := parseSpeakers(reader)
			if err != nil {
				c <- Result{Error: err}
				continue
			}

			for _, s := range speakers {
				reader, err := w.g.GetSpeaker(s.ProfilePage)
				if err != nil {
					c <- Result{Error: err}
					continue
				}

				speaker, err := parseSpeaker(reader)
				if err != nil {
					c <- Result{Error: err}
					continue
				}

				c <- Result{
					Speaker: Speaker{
						Name:         s.Name,
						Bio:          speaker.Bio,
						ProfilePage:  s.ProfilePage,
						ProfileImage: speaker.ProfileImage,
						Year:         y,
						Links:        speaker.Links,
					},
				}
			}
		}

		close(c)
	}()

	return c
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

	speaker := &Speaker{}

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

	return speaker, nil
}

func noTrimJoiner(s []string) string { return strings.Join(s, "") }

func h1Matcher(n *html.Node) bool  { return n.DataAtom == atom.H1 }
func h3Matcher(n *html.Node) bool  { return n.DataAtom == atom.H3 }
func aMatcher(n *html.Node) bool   { return n.DataAtom == atom.A }
func pMatcher(n *html.Node) bool   { return n.DataAtom == atom.P }
func imgMatcher(n *html.Node) bool { return n.DataAtom == atom.Img }
func ulMatcher(n *html.Node) bool  { return n.DataAtom == atom.Ul }

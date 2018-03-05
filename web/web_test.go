package web

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type localGetter struct {
	speakersHTMLPage string
	errSpeakers      error
	speakerHTMLPage  string
	errSpeaker       error
}

func (g *localGetter) readHTML(page string, err error) (io.Reader, error) {
	if err != nil {
		return nil, err
	}
	return os.Open(page)
}

func (g *localGetter) GetSpeakersByYear(year int) (io.Reader, error) {
	return g.readHTML(g.speakersHTMLPage, g.errSpeakers)
}

func (g *localGetter) GetSpeaker(profilePage string) (io.Reader, error) {
	return g.readHTML(g.speakerHTMLPage, g.errSpeaker)
}

func TestGetSpeakers(t *testing.T) {
	tt := []struct {
		name            string
		speakerGetter   speakerGetter
		expectedResults []Result
	}{
		{
			name: "happy path",
			speakerGetter: &localGetter{
				speakersHTMLPage: "speakers_1.htm",
				speakerHTMLPage:  "speaker.htm",
			},
			expectedResults: []Result{{}},
		},
		{
			name: "happy path",
			speakerGetter: &localGetter{
				errSpeakers: errors.New("error from speakers"),
			},
			expectedResults: []Result{{Error: errors.New("error from speakers")}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			srv := &WebSpeakerService{tc.speakerGetter}
			resultChan := srv.GetSpeakers()

			results := make([]Result, 0)
			for r := range resultChan {
				results = append(results, r)
			}
			fmt.Println("Got", len(results), "results")

			assert.Equal(t, tc.expectedResults, results)
		})
	}
}

func TestParseSpeakersPage(t *testing.T) {
	f, _ := os.Open("fosdem-speakers.htm")
	speakers, err := parseSpeakers(f)
	fmt.Println(speakers, err)
}

func TestParseSpeakerPage(t *testing.T) {
	f, _ := os.Open("speaker.htm")
	speakers, err := parseSpeaker(f)
	fmt.Println(speakers, err)
}

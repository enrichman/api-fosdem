package web

import (
	"errors"
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
			expectedResults: []Result{{
				Speaker: Speaker{
					Name:         "BSDCG Team",
					Bio:          "This is a BIO",
					ProfilePage:  "/2018/schedule/speaker/bsdcg_team/",
					ProfileImage: "/2018/schedule/speaker/francesc_campoy/cac4fd830f6d7dd839e1a8cd77ad17c9f5ba9bb39b9c2bc44b05f4568a72a1b6.jpg",
					Year:         2018,
					Links:        []Link{Link{Title: "justforfunc", URL: "http://justforfunc.com"}},
				},
			}},
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

			assert.Equal(t, tc.expectedResults, results)
		})
	}
}

package web

import (
	"fmt"
	"os"
	"testing"
)

func TestGetSpeakers(t *testing.T) {
	resultChan := GetSpeakers()
	for r := range resultChan {
		if r.Error != nil {
			fmt.Println(r.Error.Error())
		} else {
			fmt.Println(r.Speaker.Name)
		}
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

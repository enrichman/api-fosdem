package indexer_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/enrichman/api-fosdem/indexer"
)

func TestParseScheduleXML(t *testing.T) {
	speakersFile, err := os.Open("../schedule.xml")
	if err != nil {
		panic(err)
	}
	speakers := indexer.ParseScheduleXML(speakersFile)
	fmt.Println(speakers)
}

func TestParseSpeakersPage(t *testing.T) {
	speakersPage, err := os.Open("../fosdem-speakers.htm")
	if err != nil {
		panic(err)
	}
	speakersLink := indexer.ParseSpeakersPage(speakersPage)
	fmt.Println(speakersLink)
}

func TestFillSpeakersInfo(t *testing.T) {
	speakersFile, err := os.Open("../schedule.xml")
	if err != nil {
		panic(err)
	}
	speakers := indexer.ParseScheduleXML(speakersFile)

	speakersPage, err := os.Open("../fosdem-speakers.htm")
	if err != nil {
		panic(err)
	}
	speakersLink := indexer.ParseSpeakersPage(speakersPage)

	speakers = indexer.FillSpeakersInfo(speakers, speakersLink)

	fmt.Println(speakers)
}

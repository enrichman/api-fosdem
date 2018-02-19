package indexer_test

import (
	"os"
	"testing"

	"github.com/enrichman/api-fosdem/indexer"
)

func TestParseScheduleXML(t *testing.T) {
	speakersFile, err := os.Open("../schedule.xml")
	if err != nil {
		panic(err)
	}
	_ = indexer.ParseScheduleXML(speakersFile)
}

func TestParseSpeakersPage(t *testing.T) {
	speakersPage, err := os.Open("../fosdem-speakers.htm")
	if err != nil {
		panic(err)
	}
	_ = indexer.ParseSpeakersPage(speakersPage)
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

	_ = indexer.FillSpeakersInfo(speakers, speakersLink)
}

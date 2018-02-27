package web

import (
	"fmt"
	"os"
	"testing"
)

func TestParseSpeakersPage(t *testing.T) {
	f, _ := os.Open("fosdem-speakers.htm")
	speakers, err := ParseSpeakersPage(f)
	fmt.Println(speakers, err)
}

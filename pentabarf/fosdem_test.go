package pentabarf

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	f, _ := os.Open("pentabarf_test.xml")
	s, err := Parse(f)
	fmt.Printf("%+v %+v\n", s, err)
}

func Test_parseEvent(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Brussels")
	day := time.Date(2018, time.February, 3, 0, 0, 0, 0, location)
	resEv, err := parseEvent(&Event{StartStr: "12:00", DurationStr: "00:23"}, day, location)

	fmt.Printf("%+v %+v %+v\n", resEv, err, day)
}

package pentabarf

import (
	"encoding/xml"
	"io"
	"strconv"
	"strings"
	"time"
)

// Schedule maps the schedule on the XML
type Schedule struct {
	Conference *Conference `xml:"conference"`
	Days       []*Day      `xml:"day"`
}

type Conference struct {
	Title               string `xml:"title"`
	Subtitle            string `xml:"subtitle"`
	Venue               string `xml:"venue"`
	StartDate           time.Time
	StartDateStr        string `xml:"start"`
	EndDate             time.Time
	EndDateStr          string `xml:"end"`
	Days                int    `xml:"days"`
	DayChange           time.Duration
	DayChangeStr        string `xml:"day_change"`
	TimeslotDuration    time.Duration
	TimeslotDurationStr string `xml:"timeslot_duration"`
}

// Day maps the day on the XML
type Day struct {
	Index   int `xml:"index,attr"`
	Date    time.Time
	DateStr string  `xml:"date,attr"`
	Rooms   []*Room `xml:"room"`
}

// Room maps the room on the XML
type Room struct {
	Name   string   `xml:"name,attr"`
	Events []*Event `xml:"event"`
}

// Event maps the event on the XML
type Event struct {
	ID          int `xml:"id,attr"`
	Start       time.Time
	StartStr    string `xml:"start"`
	Duration    time.Duration
	DurationStr string     `xml:"duration"`
	Room        string     `xml:"room"`
	Slug        string     `xml:"slug"`
	Title       string     `xml:"title"`
	Subtitle    string     `xml:"subtitle"`
	Track       string     `xml:"track"`
	Type        string     `xml:"type"`
	Language    string     `xml:"language"`
	Abstract    string     `xml:"abstract"`
	Description string     `xml:"description"`
	Speakers    []*Speaker `xml:"persons>person"`
	Links       []*Link    `xml:"links>link"`
}

// Speaker maps the speaker on the XML and the other attributes
type Speaker struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:",chardata"`
}

type Link struct {
	URL  string `xml:"href,attr"`
	Text string `xml:",chardata"`
}

func Parse(xmlReader io.Reader) (*Schedule, error) {
	var schedule Schedule
	var err error

	err = xml.NewDecoder(xmlReader).Decode(&schedule)
	if err != nil {
		return nil, err
	}

	location, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		return nil, err
	}

	schedule.Conference, err = parseConference(schedule.Conference, location)
	if err != nil {
		return nil, err
	}

	for i, d := range schedule.Days {
		d.Date, err = time.ParseInLocation("2006-01-02", d.DateStr, location)
		if err != nil {
			return nil, err
		}
		schedule.Days[i] = d
	}

	return &schedule, nil
}

func parseConference(c *Conference, location *time.Location) (*Conference, error) {
	var err error

	c.StartDate, err = time.ParseInLocation("2006-01-02", c.StartDateStr, location)
	if err != nil {
		return nil, err
	}

	c.EndDate, err = time.ParseInLocation("2006-01-02", c.EndDateStr, location)
	if err != nil {
		return nil, err
	}

	c.DayChange, err = getDurationByString(c.DayChangeStr)
	if err != nil {
		return nil, err
	}

	c.TimeslotDuration, err = getDurationByString(c.TimeslotDurationStr)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func parseEvent(e *Event, day time.Time, location *time.Location) (*Event, error) {
	var err error

	hhMM, err := getDurationByString(e.StartStr)
	if err != nil {
		return nil, err
	}
	e.Start = day.Add(hhMM)

	e.Duration, err = getDurationByString(e.DurationStr)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func getDurationByString(str string) (time.Duration, error) {
	var dur time.Duration

	hhMMss := strings.Split(str, ":")

	hh, err := strconv.Atoi(hhMMss[0])
	if err != nil {
		return 0, err
	}
	dur += time.Hour * time.Duration(hh)

	MM, err := strconv.Atoi(hhMMss[1])
	if err != nil {
		return 0, err
	}
	dur += time.Minute * time.Duration(MM)

	if len(hhMMss) > 2 {
		ss, err := strconv.Atoi(hhMMss[2])
		if err != nil {
			return 0, err
		}
		dur += time.Second * time.Duration(ss)
	}

	return dur, nil
}

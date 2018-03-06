package pentabarf

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL        = "https://fosdem.org"
	yyyyMMddFormat = "2006-01-02"
)

type CachedScheduleService struct {
	lastModified   string
	cachedSchedule *Schedule
}

func (c *CachedScheduleService) GetSchedule(year int) (*Schedule, error) {
	url := baseURL + "/" + strconv.Itoa(year) + "/schedule/xml"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if c.lastModified != "" && c.cachedSchedule != nil {
		req.Header.Set("If-Modified-Since", c.lastModified)
	}

	scheduleResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer scheduleResp.Body.Close()

	if scheduleResp.StatusCode == http.StatusNotModified {
		return c.cachedSchedule, nil
	}

	if scheduleResp.StatusCode == http.StatusOK {
		parsedSchedule, err := Parse(scheduleResp.Body)
		if err != nil {
			return nil, err
		}

		c.lastModified = scheduleResp.Header.Get("Last-Modified")
		c.cachedSchedule = parsedSchedule

		return c.cachedSchedule, nil
	}

	return nil, errors.New("error from Fosdem server: " + strconv.Itoa(scheduleResp.StatusCode))
}

// Schedule contains the info about the Conference and the schedule of the days
type Schedule struct {
	Conference *Conference `xml:"conference"`
	Days       []*Day      `xml:"day"`

	personsMap map[string]*Person
}

// GetAllEvents returns all the events of the schedule
func (s *Schedule) GetAllEvents() []*Event {
	events := make([]*Event, 0)
	for _, d := range s.Days {
		events = append(events, d.GetAllEvents()...)
	}
	return events
}

// GetAllRooms returns all the rooms of the schedule
func (s *Schedule) GetAllRooms() []*Room {
	rooms := make([]*Room, 0)
	for _, d := range s.Days {
		rooms = append(rooms, d.Rooms...)
	}
	return rooms
}

func (s *Schedule) GetAllPersons() []*Person {
	if s.personsMap == nil {
		s.personsMap = make(map[string]*Person)
		for _, d := range s.Days {
			for _, p := range d.GetAllPersons() {
				s.personsMap[p.Name] = p
			}
		}
	}

	persons := make([]*Person, 0)
	for _, p := range s.personsMap {
		persons = append(persons, p)
	}
	return persons
}

func (s *Schedule) GetPersonByName(name string) (*Person, bool) {
	if s.personsMap == nil {
		s.personsMap = make(map[string]*Person)
		for _, d := range s.Days {
			for _, p := range d.GetAllPersons() {
				s.personsMap[p.Name] = p
			}
		}
	}
	p, ok := s.personsMap[name]
	return p, ok
}

// Conference contains the main information about the conference
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

// Day contains the info about the rooms occupied during the day
type Day struct {
	Index   int `xml:"index,attr"`
	Date    time.Time
	DateStr string  `xml:"date,attr"`
	Rooms   []*Room `xml:"room"`

	personsMap map[string]*Person
}

func (d *Day) String() string {
	return `Day{Day: "` + d.DateStr + `"}`
}

// GetAllEvents returns all the events of the day
func (d *Day) GetAllEvents() []*Event {
	events := make([]*Event, 0)
	for _, r := range d.Rooms {
		events = append(events, r.Events...)
	}
	return events
}

func (d *Day) GetAllPersons() []*Person {
	if d.personsMap == nil {
		d.personsMap = make(map[string]*Person)
		for _, r := range d.Rooms {
			for _, p := range r.GetAllPersons() {
				d.personsMap[p.Name] = p
			}
		}
	}

	persons := make([]*Person, 0)
	for _, p := range d.personsMap {
		persons = append(persons, p)
	}
	return persons
}

// Room contains all the events of the day in the room
type Room struct {
	Name   string   `xml:"name,attr"`
	Events []*Event `xml:"event"`

	personsMap map[string]*Person
}

func (r *Room) String() string {
	return `Room{Name: "` + r.Name + `"}`
}

func (r *Room) GetAllPersons() []*Person {
	if r.personsMap == nil {
		r.personsMap = make(map[string]*Person)
		for _, e := range r.Events {
			for _, p := range e.Persons {
				r.personsMap[p.Name] = p
			}
		}
	}

	persons := make([]*Person, 0)
	for _, p := range r.personsMap {
		persons = append(persons, p)
	}
	return persons
}

// Event contains all the details about the event
type Event struct {
	ID          int `xml:"id,attr"`
	Start       time.Time
	StartStr    string `xml:"start"`
	Duration    time.Duration
	DurationStr string    `xml:"duration"`
	Room        string    `xml:"room"`
	Slug        string    `xml:"slug"`
	Title       string    `xml:"title"`
	Subtitle    string    `xml:"subtitle"`
	Track       string    `xml:"track"`
	Type        string    `xml:"type"`
	Language    string    `xml:"language"`
	Abstract    string    `xml:"abstract"`
	Description string    `xml:"description"`
	Persons     []*Person `xml:"persons>person"`
	Links       []*Link   `xml:"links>link"`
}

func (e *Event) String() string {
	return `Event{Title: "` + e.Title + `"}`
}

// Person is a person of an Event
type Person struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:",chardata"`
}

// Link is a link of an Event
type Link struct {
	URL  string `xml:"href,attr"`
	Text string `xml:",chardata"`
}

// Parse will parse the Pentabarf XML returning the correspoding Schedule
// using the Europe/Brussels location (FOSDEM)
func Parse(xmlReader io.Reader) (*Schedule, error) {
	location, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		return nil, err
	}
	return ParseInLocation(xmlReader, location)
}

// ParseInLocation will parse the Pentabarf XML returning the correspoding Schedule
func ParseInLocation(xmlReader io.Reader, location *time.Location) (*Schedule, error) {
	var schedule Schedule

	err := xml.NewDecoder(xmlReader).Decode(&schedule)
	if err != nil {
		return nil, err
	}

	schedule.Conference, err = parseConference(schedule.Conference, location)
	if err != nil {
		return nil, err
	}

	for dIndex, d := range schedule.Days {
		d, err = parseDay(d, location)
		if err != nil {
			return nil, err
		}

		for rIndex, r := range d.Rooms {
			for eIndex, e := range r.Events {
				e, err = parseEvent(e, d.Date, location)
				if err != nil {
					return nil, err
				}
				r.Events[eIndex] = e
			}
			d.Rooms[rIndex] = r
		}
		schedule.Days[dIndex] = d
	}

	return &schedule, nil
}

func parseConference(c *Conference, location *time.Location) (*Conference, error) {
	var err error

	c.StartDate, err = time.ParseInLocation(yyyyMMddFormat, c.StartDateStr, location)
	if err != nil {
		return nil, err
	}

	c.EndDate, err = time.ParseInLocation(yyyyMMddFormat, c.EndDateStr, location)
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

func parseDay(d *Day, location *time.Location) (*Day, error) {
	var err error
	d.Date, err = time.ParseInLocation(yyyyMMddFormat, d.DateStr, location)
	if err != nil {
		return nil, err
	}
	return d, nil
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

	if len(hhMMss) > 1 {
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
	}

	return dur, nil
}

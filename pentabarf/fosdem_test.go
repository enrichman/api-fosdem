package pentabarf

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	f, _ := os.Open("pentabarf_test.xml")
	s, err := Parse(f)
	fmt.Printf("%+v %+v\n", s, err)
}

func Test_parseEvent(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Brussels")

	type args struct {
		*Event
		time.Time
		*time.Location
	}
	tt := []struct {
		name     string
		args     args
		expEvent *Event
		wantErr  bool
	}{
		{
			name: "happy path",
			args: args{
				Event:    &Event{StartStr: "15:30", DurationStr: "01:30"},
				Time:     time.Date(1989, time.October, 2, 0, 0, 0, 0, location),
				Location: location,
			},
			expEvent: &Event{
				Start:       time.Date(1989, time.October, 2, 15, 30, 0, 0, location),
				StartStr:    "15:30",
				Duration:    time.Duration(1)*time.Hour + time.Duration(30)*time.Minute,
				DurationStr: "01:30",
			},
		},
		{
			name: "wrong start",
			args: args{
				Event:    &Event{StartStr: "AAA"},
				Time:     time.Date(1989, time.October, 2, 0, 0, 0, 0, location),
				Location: location,
			},
			wantErr: true,
		},
		{
			name: "happy duration",
			args: args{
				Event:    &Event{StartStr: "15:30", DurationStr: "AAA"},
				Time:     time.Date(1989, time.October, 2, 0, 0, 0, 0, location),
				Location: location,
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ev, err := parseEvent(tc.args.Event, tc.args.Time, tc.args.Location)

			assert.Equal(t, tc.expEvent, ev)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_getDurationByString(t *testing.T) {
	tt := []struct {
		name        string
		args        string
		expDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "happy path full duration",
			args:        "01:12:34",
			expDuration: time.Duration(1)*time.Hour + time.Duration(12)*time.Minute + time.Duration(34)*time.Second,
		},
		{
			name:        "happy path hhMM duration",
			args:        "02:34",
			expDuration: time.Duration(2)*time.Hour + time.Duration(34)*time.Minute,
		},
		{
			name:        "happy path hh duration",
			args:        "03",
			expDuration: time.Duration(3) * time.Hour,
		},
		{
			name:    "wrong duratin string",
			args:    "asfsd",
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			duration, err := getDurationByString(tc.args)

			assert.Equal(t, tc.expDuration, duration)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

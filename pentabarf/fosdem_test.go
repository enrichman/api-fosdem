package pentabarf

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCached(t *testing.T) {
	srv := &CachedScheduleService{}
	fmt.Println(srv.GetSchedule(2018))
	fmt.Println(srv.GetSchedule(2018))
}

func Test_parseConference(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Brussels")

	type args struct {
		*Conference
		*time.Location
	}
	tt := []struct {
		name          string
		args          args
		expConference *Conference
		wantErr       bool
	}{
		{
			name: "happy path",
			args: args{
				Conference: &Conference{
					StartDateStr:        "2018-02-03",
					EndDateStr:          "2018-02-04",
					DayChangeStr:        "09:00:00",
					TimeslotDurationStr: "00:10:00",
				},
				Location: location,
			},
			expConference: &Conference{
				StartDate:           time.Date(2018, time.February, 3, 0, 0, 0, 0, location),
				StartDateStr:        "2018-02-03",
				EndDate:             time.Date(2018, time.February, 4, 0, 0, 0, 0, location),
				EndDateStr:          "2018-02-04",
				DayChange:           time.Duration(9) * time.Hour,
				DayChangeStr:        "09:00:00",
				TimeslotDuration:    time.Duration(10) * time.Minute,
				TimeslotDurationStr: "00:10:00",
			},
		},
		{
			name:    "wrong start date",
			args:    args{Conference: &Conference{StartDateStr: "AAA"}, Location: location},
			wantErr: true,
		},
		{
			name:    "wrong end date",
			args:    args{Conference: &Conference{EndDateStr: "AAA"}, Location: location},
			wantErr: true,
		},
		{
			name:    "wrong day change",
			args:    args{Conference: &Conference{DayChangeStr: "AAA"}, Location: location},
			wantErr: true,
		},
		{
			name:    "wrong timeslot duration",
			args:    args{Conference: &Conference{TimeslotDurationStr: "AAA"}, Location: location},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conf, err := parseConference(tc.args.Conference, tc.args.Location)

			assert.Equal(t, tc.expConference, conf)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_parseDay(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Brussels")

	type args struct {
		*Day
		*time.Location
	}
	tt := []struct {
		name    string
		args    args
		expDay  *Day
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{Day: &Day{DateStr: "1989-10-02"}, Location: location},
			expDay: &Day{
				Date:    time.Date(1989, time.October, 2, 0, 0, 0, 0, location),
				DateStr: "1989-10-02",
			},
		},
		{
			name:    "wrong date",
			args:    args{Day: &Day{DateStr: "aaa"}, Location: location},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			day, err := parseDay(tc.args.Day, tc.args.Location)

			assert.Equal(t, tc.expDay, day)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
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
			name: "wrong duration",
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

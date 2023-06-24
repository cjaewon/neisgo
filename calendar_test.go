package neisgo_test

import (
	"testing"
	"time"

	"github.com/cjaewon/neisgo"
)

func TestGetCalendar(t *testing.T) {
	cases := []struct {
		year  int
		month time.Month
		len   int
	}{
		{
			year:  2023,
			month: 2,
			len:   28,
		},
		{
			year:  2023,
			month: 1,
			len:   31,
		},
	}

	neis := neisgo.New(apiKey)
	neis.Set("C10", "7150144")

	for i, c := range cases {
		calendars, err := neis.GetCalendar(c.year, c.month)
		if err != nil {
			t.Fatal(err)
		}

		if len(calendars) != c.len {
			t.Fatalf("catched unexpected calendars length, case %d", i)
		}
	}
}

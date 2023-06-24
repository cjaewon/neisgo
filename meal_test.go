package neisgo_test

import (
	"testing"
	"time"

	"github.com/cjaewon/neisgo"
)

const (
	apiKey = ""
)

func TestGetMeal(t *testing.T) {
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
		meals, err := neis.GetMeal(c.year, c.month)

		if err != nil {
			t.Fatal(err)
		}

		if len(meals) != c.len {
			t.Fatalf("catched unexpected meals length, case %d", i)
		}
	}
}

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
	neis := neisgo.New(apiKey)
	neis.Set("", "")

	today := time.Now()
	tomorrow := time.Now().AddDate(0, 0, 10)

	meals, err := neis.GetMeal(today, tomorrow)
	if err != nil {
		t.Error(err)
	}

	if meals == nil {
		t.Error("failed to get meal")
	}
}

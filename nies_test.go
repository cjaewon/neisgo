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
	neis.Set("T10", "9296071")

	today := time.Now()
	tomorrow := time.Now().AddDate(0, 0, 5)

	_, err := neis.GetMeal(today, tomorrow)
	if err != nil {
		t.Error(err)
	}

}

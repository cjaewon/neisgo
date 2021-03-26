package neisgo_test

import (
	"neisgo"
	"testing"
	"time"
)

const (
	apiKey = ""
)

func TestGetMeal(t *testing.T) {
	neis := neisgo.New(apiKey)
	neis.Set("T10", "9296071")

	today := time.Now()
	tomorrow := time.Now().AddDate(0, 0, 1)

	meals, err := neis.GetMeal(today, tomorrow)
	if err != nil {
		t.Error(err)
	}

	if meals == nil {
		t.Error("failed to get meal")
	}
}

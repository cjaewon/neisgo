package neisgo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type MealTime struct {
	Breakfast string
	Lunch     string
	Dinner    string
}

type Meal struct {
	Date        time.Time
	Origin      MealTime
	Ingredients MealTime
	MealTime
}

// Text merges breakfast, lunch and dinner with a provided or default template
func (m *Meal) Text() string {
	var text string

	if m.Breakfast != "" {
		text += "[조식]\n" + m.Breakfast
		if m.Lunch != "" || m.Dinner != "" {
			text += "\n\n"
		}
	}
	if m.Lunch != "" {
		text += "[중식]\n" + m.Lunch
		if m.Dinner != "" {
			text += "\n\n"
		}
	}
	if m.Dinner != "" {
		text += "[석식]\n" + m.Dinner
	}

	return text
}

type mealSchema struct {
	Mealservicedietinfo []struct {
		Head []struct {
			ListTotalCount int `json:"list_total_count,omitempty"`
			Result         struct {
				Code    string `json:"CODE"`
				Message string `json:"MESSAGE"`
			} `json:"RESULT,omitempty"`
		} `json:"head,omitempty"`
		Row []struct {
			AtptOfcdcScCode string `json:"ATPT_OFCDC_SC_CODE"`
			AtptOfcdcScNm   string `json:"ATPT_OFCDC_SC_NM"`
			SdSchulCode     string `json:"SD_SCHUL_CODE"`
			SchulNm         string `json:"SCHUL_NM"`
			MmealScCode     string `json:"MMEAL_SC_CODE"`
			MmealScNm       string `json:"MMEAL_SC_NM"`
			MlsvYmd         string `json:"MLSV_YMD"`
			MlsvFgr         string `json:"MLSV_FGR"`
			DdishNm         string `json:"DDISH_NM"`
			OrplcInfo       string `json:"ORPLC_INFO"`
			CalInfo         string `json:"CAL_INFO"`
			NtrInfo         string `json:"NTR_INFO"`
			MlsvFromYmd     string `json:"MLSV_FROM_YMD"`
			MlsvToYmd       string `json:"MLSV_TO_YMD"`
		} `json:"row,omitempty"`
	} `json:"mealServiceDietInfo"`
}

// GetMeal gets meal data from neis
// returns Meal type array of start to end date duration length
func (n *Neis) GetMeal(year int, month time.Month) ([]Meal, error) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
	end := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Now().Location())
	duration := int(end.Sub(start).Hours()/24) + 1

	q := url.Values{
		"KEY":                []string{n.apiKey},
		"Type":               []string{"json"},
		"ATPT_OFCDC_SC_CODE": []string{n.region},
		"SD_SCHUL_CODE":      []string{n.code},
		"MLSV_FROM_YMD":      []string{start.Format("20060102")},
		"MLSV_TO_YMD":        []string{end.Format("20060102")},
	}

	url := fmt.Sprintf("https://open.neis.go.kr/hub/mealServiceDietInfo?%s", q.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data mealSchema
	var meals = make([]Meal, duration)

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	for _, row := range data.Mealservicedietinfo[1].Row {
		d, err := time.Parse("20060102", row.MlsvYmd)
		if err != nil {
			return nil, err
		}

		first := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
		index := (d.Day() - start.Day()) + (int(d.Month())-int(start.Month()))*first.AddDate(0, 1, -1).Day()

		if reflect.ValueOf(meals[index]).IsZero() {
			meals[index] = Meal{
				Date: d,
			}
		}

		switch row.MmealScCode {
		case "1":
			meals[index].Breakfast = genPlainText(row.DdishNm)
			meals[index].Origin.Breakfast = genPlainText(row.OrplcInfo)
			meals[index].Ingredients.Breakfast = genPlainText(row.NtrInfo)
		case "2":
			meals[index].Lunch = genPlainText(row.DdishNm)
			meals[index].Origin.Lunch = genPlainText(row.OrplcInfo)
			meals[index].Ingredients.Lunch = genPlainText(row.NtrInfo)
		case "3":
			meals[index].Dinner = genPlainText(row.DdishNm)
			meals[index].Origin.Dinner = genPlainText(row.OrplcInfo)
			meals[index].Ingredients.Dinner = genPlainText(row.NtrInfo)
		}
	}

	return meals, nil
}

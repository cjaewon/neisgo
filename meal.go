package neisgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const (
	mealTmpl = `
{{ if .Breakfast }}
[조식]
{{ .Breakfast }}
{{ end }}
{{ if .Lunch }}
[중식]
{{ .Lunch }}
{{ end }}
{{ if .Dinner }}
[석식]
{{ .Dinner }}
{{ end }}`
)

type MealTime struct {
	Breakfast string
	Lunch     string
	Dinner    string
}

type Meal struct {
	EducationCenter string
	SchoolName      string
	Date            time.Time
	Origin          MealTime
	Ingredients     MealTime
	MealTime
}

func (m *Meal) Text() (string, error) {
	tmpl := template.New("meal")

	tmpl, err := tmpl.Parse(mealTmpl)
	if err != nil {
		return "", err
	}

	b := bytes.NewBufferString("")
	if err := tmpl.Execute(b, m); err != nil {
		return "", err
	}

	return b.String(), nil
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
func (n *Neis) GetMeal(start, end time.Time) ([]Meal, error) {
	queryParams := url.Values{
		"KEY":                []string{n.apiKey},
		"Type":               []string{"json"},
		"ATPT_OFCDC_SC_CODE": []string{n.region},
		"SD_SCHUL_CODE":      []string{n.code},
		"MLSV_FROM_YMD":      []string{start.Format("20060102")},
		"MLSV_TO_YMD":        []string{end.Format("20060102")},
	}

	url := fmt.Sprintf("https://open.neis.go.kr/hub/mealServiceDietInfo?%s", queryParams.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	duration := int(end.Sub(start).Hours() / 24)

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
				EducationCenter: row.AtptOfcdcScCode,
				SchoolName:      row.SchulNm,
				Date:            d,
			}
		}

		switch row.MmealScCode {
		case "1":
			meals[index].Breakfast = row.DdishNm
			meals[index].Origin.Breakfast = row.OrplcInfo
			meals[index].Ingredients.Breakfast = row.NtrInfo
		case "2":
			meals[index].Lunch = row.DdishNm
			meals[index].Origin.Lunch = row.OrplcInfo
			meals[index].Ingredients.Lunch = row.NtrInfo
		case "3":
			meals[index].Dinner = row.DdishNm
			meals[index].Origin.Dinner = row.OrplcInfo
			meals[index].Ingredients.Dinner = row.NtrInfo
		}
	}

	return meals, nil
}

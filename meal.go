package neisgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Meal struct {
	EducationCenter string
	SchoolName      string

	Date        time.Time
	Text        string
	Type        string
	Origin      string
	Ingredients string
}

type MealSchema struct {
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

// GetMeal getes meal data from neis
func (n *Neis) GetMeal(t ...time.Time) (*[]Meal, error) {
	if len(t) <= 0 {
		return nil, errors.New("missing time parameters")
	}

	queryParams := url.Values{
		"KEY":                []string{n.apiKey},
		"Type":               []string{"json"},
		"ATPT_OFCDC_SC_CODE": []string{n.region},
		"SD_SCHUL_CODE":      []string{n.code},
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

	var data MealSchema
	var meals []Meal

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	for _, row := range data.Mealservicedietinfo[1].Row {
		d, err := time.Parse("20060102", row.MlsvYmd)
		if err != nil {
			return nil, err
		}

		meals = append(meals, Meal{
			EducationCenter: row.AtptOfcdcScCode,
			SchoolName:      row.SchulNm,
			Date:            d,
			Text:            row.DdishNm,
			Type:            row.MmealScNm,
			Origin:          row.OrplcInfo,
			Ingredients:     row.NtrInfo,
		})
	}

	return &meals, nil
}

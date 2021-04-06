package neisgo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Calendar represents calendar of the day
type Calendar struct {
	Date time.Time
	Name string

	// Content is usually blank
	Content string

	// ClassTime represents night class or day class
	ClassTime string

	// Deduction represents class deduction
	Deduction string

	// Target is targeting school year of calendar
	// if high school or middle school, then max length will 3 else 6
	Target [6]bool
}

// IsZero returns true if c is the zero value
func (c Calendar) IsZero() bool {
	return c == (Calendar{})
}

// Text returns merged  name, content
func (c Calendar) Text() string {
	if c.Content == "" {
		return fmt.Sprintf("[%s]", c.Name)
	}
	return fmt.Sprintf("[%s]\n%s", c.Name, c.Content)
}

type calendarSchema struct {
	Schoolschedule []struct {
		Head []struct {
			ListTotalCount int `json:"list_total_count,omitempty"`
			Result         struct {
				Code    string `json:"CODE"`
				Message string `json:"MESSAGE"`
			} `json:"RESULT,omitempty"`
		} `json:"head,omitempty"`
		Row []struct {
			AtptOfcdcScCode   string `json:"ATPT_OFCDC_SC_CODE"`
			AtptOfcdcScNm     string `json:"ATPT_OFCDC_SC_NM"`
			SdSchulCode       string `json:"SD_SCHUL_CODE"`
			SchulNm           string `json:"SCHUL_NM"`
			Ay                string `json:"AY"`
			DghtCrseScNm      string `json:"DGHT_CRSE_SC_NM"`
			SchulCrseScNm     string `json:"SCHUL_CRSE_SC_NM"`
			SbtrDdScNm        string `json:"SBTR_DD_SC_NM"`
			AaYmd             string `json:"AA_YMD"`
			EventNm           string `json:"EVENT_NM"`
			EventCntnt        string `json:"EVENT_CNTNT"`
			OneGradeEventYn   string `json:"ONE_GRADE_EVENT_YN"`
			TwGradeEventYn    string `json:"TW_GRADE_EVENT_YN"`
			ThreeGradeEventYn string `json:"THREE_GRADE_EVENT_YN"`
			FrGradeEventYn    string `json:"FR_GRADE_EVENT_YN"`
			FivGradeEventYn   string `json:"FIV_GRADE_EVENT_YN"`
			SixGradeEventYn   string `json:"SIX_GRADE_EVENT_YN"`
			LoadDtm           string `json:"LOAD_DTM"`
		} `json:"row,omitempty"`
	} `json:"SchoolSchedule"`
}

// GetCalendar gets calendar data from neis
func (n *Neis) GetCalendar(year int, month time.Month) ([]Calendar, error) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
	end := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Now().Location())
	duration := int(end.Sub(start).Hours()/24) + 1

	q := url.Values{
		"KEY":                []string{n.apiKey},
		"Type":               []string{"json"},
		"ATPT_OFCDC_SC_CODE": []string{n.region},
		"SD_SCHUL_CODE":      []string{n.code},
		"AA_FROM_YMD":        []string{start.Format("20060102")},
		"AA_TO_YMD":          []string{end.Format("20060102")},
	}

	url := fmt.Sprintf("https://open.neis.go.kr/hub/SchoolSchedule?%s", q.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data calendarSchema
	var calendars = make([]Calendar, duration)

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	for _, row := range data.Schoolschedule[1].Row {
		d, err := time.Parse("20060102", row.AaYmd)
		if err != nil {
			return nil, err
		}

		index := end.Day() - d.Day()

		calendars[index] = Calendar{
			Date:      d,
			Name:      row.EventNm,
			Content:   row.EventCntnt,
			ClassTime: row.DghtCrseScNm,
			Deduction: row.SbtrDdScNm,
			Target: [6]bool{
				row.OneGradeEventYn == "Y",
				row.TwGradeEventYn == "Y",
				row.ThreeGradeEventYn == "Y",
				row.FrGradeEventYn == "Y",
				row.FivGradeEventYn == "Y",
				row.SixGradeEventYn == "Y",
			},
		}

	}

	return calendars, nil
}

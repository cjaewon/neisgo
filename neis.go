package neisgo

type Neis struct {
	apiKey string

	code   string
	region string

	MealTmpl     string
	CalendarTmpl string
}

// New creates a Neis instance
func New(apiKey string) *Neis {
	n := Neis{
		apiKey: apiKey,
	}

	n.MealTmpl = `
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
{{ end }}
`

	return &n
}

// Set sets a school which will use
func (n *Neis) Set(region, code string) {
	n.code = code
	n.region = region
}

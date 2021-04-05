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

	return &n
}

// Set sets a school which will use
func (n *Neis) Set(region, code string) {
	n.code = code
	n.region = region
}

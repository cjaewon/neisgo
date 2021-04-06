package neisgo

// Neis handles open neis api related
type Neis struct {
	apiKey string

	code   string
	region string
}

// New creates a new Neis instance
func New(apiKey string) *Neis {
	n := Neis{
		apiKey: apiKey,
	}

	return &n
}

// Set sets a school which will use
func (n *Neis) Set(region, code string) {
	n.region = region
	n.code = code
}

package neisgo

import "errors"

type Neis struct {
	apiKey string

	code       string
	region     string
	schoolCode string
}

// New creates a Neis instance
func New(apiKey string) *Neis {
	n := Neis{
		apiKey: apiKey,
	}

	return &n
}

// Set sets a school which will use
func (n *Neis) SetSchool(code string) error {
	if len(code) < 10 {
		return errors.New("invalid code was provided")
	}

	n.code = code
	n.region = code[0:3]
	n.schoolCode = code[3:]

	return nil
}

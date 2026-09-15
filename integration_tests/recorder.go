package integration

import (
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func NewRecorderWithAuthFilter(path string) (*recorder.Recorder, error) {
	r, err := recorder.New(path)
	if err != nil {
		return nil, err
	}

	// Add a filter which removes Authorization headers from the recorded request.
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		return nil
	})

	return r, nil
}

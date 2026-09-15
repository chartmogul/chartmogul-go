package integration

import (
	"os"
	"strings"

	cm "github.com/chartmogul/chartmogul-go/v5"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

// CHARTMOGUL_API_URL points the SDK at a non-production API while recording cassettes,
// e.g. http://localhost:9292/v1. Unset in CI so replays match the recorded production URLs.
func init() {
	if apiURL := os.Getenv("CHARTMOGUL_API_URL"); apiURL != "" {
		cm.SetURL(strings.TrimSuffix(apiURL, "/") + "/%v")
	}
}

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

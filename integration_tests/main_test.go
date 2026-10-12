package integration

import (
	"os"
	"strings"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v5"
)

// CHARTMOGUL_API_URL points the SDK at a non-production API while recording cassettes,
// e.g. http://localhost:9292/v1. Unset in CI so replays match the recorded production URLs.
func TestMain(m *testing.M) {
	if apiURL := os.Getenv("CHARTMOGUL_API_URL"); apiURL != "" {
		cm.SetURL(strings.TrimSuffix(apiURL, "/") + "/%v")
	}
	os.Exit(m.Run())
}

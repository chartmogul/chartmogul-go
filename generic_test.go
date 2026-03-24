package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRetryDisabled(t *testing.T) {
	var requestCount int
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "rate limited"}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	enabled := false
	tested := &API{
		ApiKey: "token",
		Retry:  &RetryConfig{Enabled: &enabled},
	}
	err := tested.delete("path1/:uuid", "uuid1")
	if err == nil {
		t.Fatal("Expected to fail with 429")
	}
	if requestCount != 1 {
		t.Errorf("Expected exactly 1 request (no retries), got: %v", requestCount)
	}
}

func TestRetryStatusTooManyRequests(t *testing.T) {
	var i int
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if i == 0 {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte("{errors: \"nooo\"}")) //nolint
				} else if i == 1 {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("{}")) //nolint
				}
				i++
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	err := tested.delete("path1/:uuid", "uuid1")
	if err != nil {
		spew.Dump(err)
		t.Fatal("Expected to retry")
	}
}

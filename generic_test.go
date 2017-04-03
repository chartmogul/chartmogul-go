package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRetryStatusTooManyRequests(t *testing.T) {
	var i int
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if i == 0 {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte("{errors: \"nooo\"}"))
				} else if i == 1 {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("{}"))
				}
				i++
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	err := tested.delete("path1/:uuid", "uuid1")
	if err != nil {
		spew.Dump(err)
		t.Fatal("Expected to retry")
	}
}

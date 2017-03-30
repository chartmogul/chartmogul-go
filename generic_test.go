package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// do nothing
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	tested.delete("path1/:uuid", "uuid1")
}

//TODO: write unit tests for multiple parallel backed off requests

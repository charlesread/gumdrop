package internal_test

import (
	"fmt"
	"github.com/charlesread/gumdrop/internal"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(&internal.Server{})
}

func TestServer(t *testing.T) {

	server := createServer(t)
	defer server.Close()

	t.Run("Posting to /drop returns 200", func(t *testing.T) {

		response, _ := http.Post(fmt.Sprintf("%s/drop", server.URL), "multipart/form-data", nil)
		if response.StatusCode != 200 {
			t.Errorf("got %v, want %v", response.StatusCode, 200)
		}

	})

	t.Run("Posting to /dropz returns 404", func(t *testing.T) {

		response, _ := http.Post(fmt.Sprintf("%s/dropz", server.URL), "", nil)
		if response.StatusCode != 404 {
			t.Errorf("got %q, want %q", response.StatusCode, 404)
		}

	})

	t.Run("Getting /drop returns 405", func(t *testing.T) {

		response, _ := http.Get(fmt.Sprintf("%s/drop", server.URL))
		if response.StatusCode != 405 {
			t.Errorf("got %q, want %q", response.StatusCode, 405)
		}

	})

	t.Run("Getting /dropz returns 404", func(t *testing.T) {

		response, _ := http.Get(fmt.Sprintf("%s/dropz", server.URL))
		if response.StatusCode != 404 {
			t.Errorf("got %q, want %q", response.StatusCode, 404)
		}

	})

	t.Run("Content-type must be multipart/form-data", func(t *testing.T) {

		response, _ := http.Post(fmt.Sprintf("%s/drop", server.URL), "NOTmultipart/form-data", nil)
		if response.StatusCode != http.StatusUnsupportedMediaType {
			t.Errorf("got %v, want %v", response.StatusCode, http.StatusUnsupportedMediaType)
		}

	})

}

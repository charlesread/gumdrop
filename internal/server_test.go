package internal_test

import (
	"github.com/charlesread/gumdrop/internal"
	"net/http/httptest"
	"testing"
)

func createServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(&internal.Server{})
}

func TestServer(t *testing.T) {

	t.Run("server starts", func(t *testing.T) {

		server := createServer(t)

		if server == nil {
			t.Errorf("Server is nil")
		}

	})

}

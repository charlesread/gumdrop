// +build integration

package integration_test

import (
	"net/http"
	"testing"
)

func assert(t *testing.T, got interface{}, want interface{}) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIntegration(t *testing.T) {

	baseUrl := "http://localhost:8080"

	t.Run("GET / should return 405", func(t *testing.T) {
		req, _ := http.Get(baseUrl)
		got := req.StatusCode
		want := http.StatusMethodNotAllowed
		assert(t, got, want)
	})

	t.Run("POST / should return 401", func(t *testing.T) {
		req, _ := http.Post(baseUrl, "", nil)
		got := req.StatusCode
		want := http.StatusUnauthorized
		assert(t, got, want)
	})

	t.Run("POST / with wrong Authorization should return 401", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, baseUrl, nil)
		req.Header.Set("Authorization", "bearer wrong")
		client := http.Client{}
		resp, _ := client.Do(req)
		got := resp.StatusCode
		want := http.StatusUnauthorized
		assert(t, got, want)
	})

	t.Run("POST / with correct Authorization should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, baseUrl, nil)
		req.Header.Set("Authorization", "bearer someOtherEquallySuperSecretToken")
		client := http.Client{}
		resp, _ := client.Do(req)
		got := resp.StatusCode
		want := http.StatusBadRequest
		assert(t, got, want)
	})

	t.Run("POST / with correct Authorization and wrong content type should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, baseUrl, nil)
		req.Header.Set("Authorization", "bearer someOtherEquallySuperSecretToken")
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{}
		resp, _ := client.Do(req)
		got := resp.StatusCode
		want := http.StatusBadRequest
		assert(t, got, want)
	})

	t.Run("POST / with correct Authorization and correct content type should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, baseUrl, nil)
		req.Header.Set("Authorization", "bearer someOtherEquallySuperSecretToken")
		req.Header.Set("Content-Type", "multipart/form/data")
		client := http.Client{}
		resp, _ := client.Do(req)
		got := resp.StatusCode
		want := http.StatusBadRequest
		assert(t, got, want)
	})

}

package internal

import (
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewProcessResult(t *testing.T) {

	got := newProcessResult()

	want := processResult{
		err:        nil,
		msg:        "",
		success:    true,
		statusCode: http.StatusCreated,
		baseDir:    "",
		directory:  "",
	}

	if reflect.DeepEqual(*got, want) != true {
		t.Errorf("got %v, want %v", got, want)
	}

}

//r.Header.Set("Content-Type", "application/json")

func TestRequestIsValid(t *testing.T) {

	t.Run("only requests to / don't 404", func(t *testing.T) {
		pr := newProcessResult()
		r := httptest.NewRequest("POST", "/a", nil)

		want := http.StatusNotFound
		pr.requestIsValid(r)
		if pr.statusCode != want {
			t.Errorf("got %v, want %v", pr.statusCode, want)
		}

		r = httptest.NewRequest("GET", "/a", nil)
		pr.requestIsValid(r)
		if pr.statusCode != want {
			t.Errorf("got %v, want %v", pr.statusCode, want)
		}

		r = httptest.NewRequest("POST", "/", nil)
		pr = newProcessResult()
		pr.requestIsValid(r)
		if pr.statusCode == 404 {
			t.Errorf("got %v, want !404", pr.statusCode)
		}

	})

	t.Run("no Authorization header should return 401", func(t *testing.T) {
		pr := newProcessResult()
		r := httptest.NewRequest("POST", "/", nil)

		want := http.StatusUnauthorized
		pr.requestIsValid(r)
		if pr.statusCode != want {
			t.Errorf("got %v, want %v", pr.statusCode, want)
		}
	})

	t.Run("malformed Authorization header should return 401", func(t *testing.T) {
		pr := newProcessResult()
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Authorization", "nope")

		want := http.StatusUnauthorized
		pr.requestIsValid(r)
		if pr.statusCode != want {
			t.Errorf("got %v, want %v", pr.statusCode, want)
		}
	})

	t.Run("wrong token should return 401", func(t *testing.T) {
		pr := newProcessResult()
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Authorization", "bearer WRONGTOKEN")

		want := http.StatusUnauthorized
		pr.requestIsValid(r)
		if pr.statusCode != want {
			t.Errorf("got %v, want %v", pr.statusCode, want)
		}
	})

	t.Run("correct token, but nothing else, should return 400 and MSG_NO_MULTIPART", func(t *testing.T) {
		InitViper()
		viper.Set("LogFilePath", "/dev/null")
		InitLogger()
		pr := newProcessResult()
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Authorization", "bearer superSecretToken")

		want := http.StatusBadRequest
		pr.requestIsValid(r)
		if pr.statusCode != want {
			t.Errorf("got %v, want %v", pr.statusCode, want)
		}
		if pr.msg != MsgNoMultipart {
			t.Errorf("got %v, want %v", pr.msg, MsgNoMultipart)
		}
	})

	t.Run("correct token, no x-directory header, 400 and MSG_NO_DIRECTORY", func(t *testing.T) {
		InitViper()
		viper.Set("LogFilePath", "/dev/null")
		InitLogger()
		pr := newProcessResult()
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Authorization", "bearer superSecretToken")
		r.Header.Set("Content-Type", "multipart/form-data")

		want := http.StatusBadRequest
		pr.requestIsValid(r)
		if pr.statusCode != want {
			t.Errorf("got %v, want %v", pr.statusCode, want)
		}
		if pr.msg != MsgNoDirectory {
			t.Errorf("got %v, want %v", pr.msg, MsgNoDirectory)
		}
	})

	t.Run("no file 400 and MSG_NO_DIRECTORY", func(t *testing.T) {
		InitViper()
		viper.Set("LogFilePath", "/dev/null")
		InitLogger()
		pr := newProcessResult()
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Authorization", "bearer superSecretToken")
		r.Header.Set("Content-Type", "multipart/form-data")
		r.Header.Set("x-directory", "testdir")

		want := http.StatusBadRequest
		pr.requestIsValid(r)
		pr.saveFiles(r)
		if pr.statusCode != want {
			t.Errorf("got %v, want %v", pr.statusCode, want)
		}
		if pr.msg != MsgNoFile {
			t.Errorf("got %v, want %v", pr.msg, MsgNoFile)
		}
	})

}

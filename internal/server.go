package internal

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
)

// generalized abstraction of the server, implements the ServeHTTP method that actually handles the request and
// response, you'll see that this method simply calls a few private helper methods
type Server struct {
}

// main request "handler"
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("Request received from %v for endpoint %v\n", r.RemoteAddr, r.URL)

	pr := newProcessResult()

	requestIsValid(pr, r)
	saveFile(pr, r)

	log.Printf("Process Result: %v\n", *pr)

	writeProcessRequest(pr, w)

}

// struct to hold results of different processing stages
type processResult struct {
	err        error
	msg        string
	success    bool
	statusCode int
	baseDir    string
	directory  string
}

// constructor function to make a "good" result by default, will be altered along the way and used to create
// the ultimate response to the client
func newProcessResult() *processResult {
	pr := &processResult{
		err:        nil,
		msg:        "",
		success:    true,
		statusCode: http.StatusOK,
		baseDir:    viper.GetString("BaseDir"),
		directory:  "",
	}
	return pr
}

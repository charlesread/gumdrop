package internal

import (
	"github.com/spf13/viper"
	"net/http"
)

// generalized abstraction of the server, implements the ServeHTTP method that actually handles the request and
// response, you'll see that this method simply calls a few private helper methods
type Server struct {
}

// main request "handler"
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	Log.Printf("Request received from %v for endpoint %v %v\n", r.RemoteAddr, r.Method, r.URL)

	pr := newProcessResult()
	pr.baseDir = viper.GetString("BaseDir")

	pr.requestIsValid(r)
	pr.saveFiles(r)

	Log.Printf("Process Result: %v\n", *pr)

	pr.writeProcessRequest(w)

}

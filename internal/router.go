package internal

import (
	"github.com/charlesread/gumdrop/internal/routes"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

var r *mux.Router
var once sync.Once

func init() {
	once = sync.Once{}
}

func NewRouter() *mux.Router {

	once.Do(func() {
		r = mux.NewRouter()
		r.Use(middlewareAuthentication)
		r.HandleFunc("/", routes.UploadHandler).Methods("POST")
		r.HandleFunc("/stat", routes.StatHandler).Methods("GET")
	})

	return r

}

func middlewareAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearerHeader := r.Header.Get("Authorization")
		if bearerHeader == "" {
			http.Error(w, MsgNoAuth, http.StatusUnauthorized)
			return
		}
		if len(bearerHeader) < 8 {
			http.Error(w, MsgAuthMalformed, http.StatusUnauthorized)
			return
		}
		token := bearerHeader[7:]
		tokenValid := validateToken(token)
		if tokenValid != true {
			http.Error(w, MsgTokenInvalid, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

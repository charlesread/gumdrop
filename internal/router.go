package internal

import (
	"github.com/charlesread/gumdrop/internal/routes"
	"github.com/gorilla/mux"
	"sync"
)

var r *mux.Router

func NewRouter() *mux.Router {

	once := sync.Once{}

	once.Do(func() {
		r = mux.NewRouter()
		r.HandleFunc("/", routes.UploadHandler).Methods("POST")
		r.HandleFunc("/stat", routes.StatHandler).Methods("GET")
	})

	return r

}

func GetRouter() *mux.Router {
	return NewRouter()
}

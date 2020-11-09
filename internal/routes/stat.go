package routes

import "net/http"

func StatHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

package routes

import "net/http"

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

package internal

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func setProcessResult(
	pr *processResult,
	err error,
	msg string,
	success bool,
	statusCode int) {
	pr.err = err
	pr.msg = msg
	pr.success = success
	pr.statusCode = statusCode
}

// basic request validation, makes sure they're requesting the right endpoint, have the right method,
// are authentic, et cetera, sure there are easier ways to do this, but this is a simple need
func requestIsValid(pr *processResult, r *http.Request) {

	if pr.success != true {
		return
	}

	// only the /drop path exists
	if r.URL.String() != "/drop" {
		setProcessResult(pr, nil, "Only the GET /drop endpoint exists.", false, http.StatusNotFound)
		return
	}

	// only POST is allowed
	if r.Method != http.MethodPost {
		setProcessResult(pr, nil, "Only the GET /drop endpoint exists.", false, http.StatusMethodNotAllowed)
		return
	}

	// make sure Authorization header is present and valid
	bearerHeader := r.Header.Get("Authorization")
	if bearerHeader == "" {
		setProcessResult(pr, nil, "No Authorization header set", false, http.StatusUnauthorized)
		return
	}
	if len(bearerHeader) < 8 {
		setProcessResult(pr, nil, "Authorization header is malformed.", false, http.StatusUnauthorized)
		return
	}
	token := bearerHeader[7:]
	tokenValid := validateToken(token)
	if tokenValid != true {
		setProcessResult(pr, nil, "Invalid token.", false, http.StatusUnauthorized)
		return
	}

	// get Content-Type header, returns "" if header does not exist
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		setProcessResult(pr, nil, "Could not detect multipart/form-data Content-Type.", false, http.StatusBadRequest)
		return
	}

	// ensure that x-directory header exists
	pr.directory = r.Header.Get("x-directory")
	if pr.directory == "" {
		setProcessResult(pr, nil, "x-directory header not sent.", false, http.StatusUnauthorized)
		return
	}

}

func validateToken(token string) bool {
	return true
}

func writeFileHeader(h *multipart.FileHeader, pr *processResult) {

	// open temporaryFile
	temporaryFile, err := h.Open()
	defer func() {
		err := temporaryFile.Close()
		log.Printf("temporaryFile closed, err: %v\n", err)
	}()
	if err != nil {
		setProcessResult(pr, err, err.Error(), false, http.StatusInternalServerError)
		return
	}

	// create directory where temporaryFile should be persisted
	dir := filepath.Join(viper.GetString("BaseDir"), pr.directory)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Printf("Error creating directory %q: %v\n", dir, err.Error())
		setProcessResult(pr, err, err.Error(), false, http.StatusInternalServerError)
		return
	}

	// actually write the temporaryFile to the directory
	writePath := filepath.Join(viper.GetString("BaseDir"), pr.directory, h.Filename)
	log.Printf("writePath: %v", writePath)
	file, err := os.Create(writePath)
	defer func() {
		err := file.Close()
		log.Printf("file closed, err: %v\n", err)
	}()
	_, err = io.Copy(file, temporaryFile)
	if err != nil {
		msg := fmt.Sprintf("Failed to Copy file: %v", err.Error())
		setProcessResult(pr, err, msg, false, http.StatusInternalServerError)
		return
	}
}

// actually persist the file to the OS
func saveFile(pr *processResult, r *http.Request) {

	if pr.success != true {
		return
	}

	err := r.ParseMultipartForm(1 << 62)

	// if there was a problem parsing the data let's just stop
	if err != nil {
		setProcessResult(pr, err, fmt.Sprintf("Failed to parse multipart message: %v", err.Error()), false, http.StatusBadRequest)
		return
	}

	// loop through all of the files (there can be more than one!) and save them permanently
	for _, h := range r.MultipartForm.File["file"] {
		writeFileHeader(h, pr)
	}

}

// basic marshaller for a processResult, just "makes json"
func createBody(pr *processResult) string {
	body := fmt.Sprintf("{\"message\":%q, \"success\": %v, \"statusCode\": %v}", pr.msg, pr.success, pr.statusCode)
	return body
}

// actually write the response to the client
func writeProcessRequest(pr *processResult, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(pr.statusCode)
	body := createBody(pr)
	_, err := w.Write([]byte(body))
	if err != nil {
		log.Printf("Error writing response: %v\n", err.Error())
	}
}

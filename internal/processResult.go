package internal

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

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
		statusCode: http.StatusCreated,
		baseDir:    "",
		directory:  "",
	}
	return pr
}

const (
	MSG_NO_ROUTE       = "Only the POST / endpoint exists."
	MSG_NO_MULTIPART   = "Could not detect multipart/form-data Content-Type."
	MSG_NO_AUTH        = "No Authorization header set"
	MSG_AUTH_MALFORMED = "Authorization header is malformed."
	MSG_TOKEN_INVALID  = "Invalid token."
	MSG_NO_DIRECTORY   = "x-directory header not sent."
	MSG_NO_FILE        = "Request must have a `file` or mulriple `file`s."
)

// utility function to set values
func (pr *processResult) setProcessResult(
	err error,
	msg string,
	success bool,
	statusCode int) {
	pr.err = err
	pr.msg = msg
	pr.success = success
	pr.statusCode = statusCode
}

func validateToken(token string) bool {
	for _, v := range viper.GetStringSlice("Tokens") {
		if v == token {
			return true
		}
	}
	return false
}

//---------

// basic request validation, makes sure they're requesting the right endpoint, have the right method,
// are authentic, et cetera, sure there are easier ways to do this, but this is a simple need
func (pr *processResult) requestIsValid(r *http.Request) {

	if pr.success != true {
		return
	}

	// only the /drop path exists
	if r.URL.String() != "/" {
		pr.setProcessResult(nil, MSG_NO_ROUTE, false, http.StatusNotFound)
		return
	}

	// only POST is allowed
	if r.Method != http.MethodPost {
		pr.setProcessResult(nil, MSG_NO_ROUTE, false, http.StatusMethodNotAllowed)
		return
	}

	// make sure Authorization header is present and valid
	bearerHeader := r.Header.Get("Authorization")
	if bearerHeader == "" {
		pr.setProcessResult(nil, MSG_NO_AUTH, false, http.StatusUnauthorized)
		return
	}
	if len(bearerHeader) < 8 {
		pr.setProcessResult(nil, MSG_AUTH_MALFORMED, false, http.StatusUnauthorized)
		return
	}
	token := bearerHeader[7:]
	tokenValid := validateToken(token)
	if tokenValid != true {
		pr.setProcessResult(nil, MSG_TOKEN_INVALID, false, http.StatusUnauthorized)
		return
	}

	// get Content-Type header, returns "" if header does not exist
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		pr.setProcessResult(nil, MSG_NO_MULTIPART, false, http.StatusBadRequest)
		return
	}

	// ensure that x-directory header exists
	pr.directory = r.Header.Get("x-directory")
	if pr.directory == "" {
		pr.setProcessResult(nil, MSG_NO_DIRECTORY, false, http.StatusBadRequest)
		return
	}

}

//---------

// actually persist the file to the OS
func (pr *processResult) saveFiles(r *http.Request) {

	if pr.success != true {
		return
	}

	err := r.ParseMultipartForm(32 << 20)

	// if there was a problem parsing the data let's just stop
	if err != nil {
		Log.Printf("Failed to parse multipart message: %v", err.Error())
		pr.setProcessResult(err, MSG_NO_FILE, false, http.StatusBadRequest)
		return
	}

	// loop through all of the files (there can be more than one!) and save them permanently
	for _, h := range r.MultipartForm.File["file"] {
		pr.writeFileHeader(h)
	}

}

func (pr *processResult) writeFileHeader(h *multipart.FileHeader) {

	// open temporaryFile
	temporaryFile, err := h.Open()
	defer func() {
		err := temporaryFile.Close()
		Log.Printf("temporaryFile closed, err: %v\n", err)
	}()
	if err != nil {
		pr.setProcessResult(err, err.Error(), false, http.StatusInternalServerError)
		return
	}

	// create directory where temporaryFile should be persisted
	dir := filepath.Join(viper.GetString("BaseDir"), pr.directory)
	err = os.MkdirAll(dir, os.FileMode(viper.GetUint32("PathMode")))
	if err != nil {
		Log.Printf("Error creating directory %q: %v\n", dir, err.Error())
		pr.setProcessResult(err, err.Error(), false, http.StatusInternalServerError)
		return
	}

	// actually write the temporaryFile to the directory
	writePath := filepath.Join(viper.GetString("BaseDir"), pr.directory, h.Filename)
	Log.Printf("writePath: %v", writePath)
	file, err := os.Create(writePath)
	_ = file.Chmod(os.FileMode(viper.GetUint32("FileMode")))
	defer func() {
		err := file.Close()
		Log.Printf("file closed, err: %v\n", err)
	}()
	_, err = io.Copy(file, temporaryFile)
	if err != nil {
		msg := fmt.Sprintf("Failed to Copy file: %v", err.Error())
		pr.setProcessResult(err, msg, false, http.StatusInternalServerError)
		return
	}

	// append message to msg indicating success in persisting file
	if len(pr.msg) != 0 {
		pr.msg += ", "
	}
	pr.msg += fmt.Sprintf("%v persisted", writePath)
}

//---------

// actually write the response to the client
func (pr *processResult) writeProcessRequest(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(pr.statusCode)
	body := pr.createBody()
	_, err := w.Write([]byte(body))
	if err != nil {
		Log.Printf("Error writing response: %v\n", err.Error())
	}
}

// basic marshaller for a processResult, just "makes json"
func (pr *processResult) createBody() string {
	body := fmt.Sprintf("{\"message\":%q, \"success\": %v, \"statusCode\": %v}", pr.msg, pr.success, pr.statusCode)
	return body
}

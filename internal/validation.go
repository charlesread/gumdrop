package internal

import "net/http"

func Validate(r *http.Request) error {

	// only the /drop path exists
	if r.URL.String() != "/" {
		return ErrNoRoute
	}

	// only POST is allowed
	if r.Method != http.MethodPost {
		return ErrNoRoute
	}

	// make sure Authorization header is present and valid
	bearerHeader := r.Header.Get("Authorization")
	if bearerHeader == "" {
		return ErrNoAuth
	}
	if len(bearerHeader) < 8 {
		return ErrAuthMalformed
	}
	token := bearerHeader[7:]
	tokenValid := validateToken(token)
	if tokenValid != true {
		return ErrTokenInvalid
	}

	// get Content-Type header, returns "" if header does not exist
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return ErrNoMultipart
	}

	// ensure that x-directory header exists
	if r.Header.Get("x-directory") == "" {
		return ErrNoDirectory
	}

	return nil

}

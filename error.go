package graceful

import (
	"net/http"
	"net/url"
)

func Error(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

func ErrorFromURL(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*url.Error); ok {
		return e.Err
	}
	return err
}

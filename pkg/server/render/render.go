package render

import (
	"fmt"
	"net/http"
)

// Header keys for http responses
const (
	ContentType = "Content-Type"
)

// Content types for plain text responses
const (
	plainContentType = "text/plain; charset=utf-8"
)

type Renderer func(code int, w http.ResponseWriter, obj any) error

func Text(code int, w http.ResponseWriter, text string) error {
	w.Header().Set(ContentType, plainContentType)
	w.WriteHeader(code)

	fmt.Fprintln(w, text)
	return nil
}

func Textf(code int, w http.ResponseWriter, text string, a ...any) error {
	w.Header().Set(ContentType, plainContentType)
	w.WriteHeader(code)

	fmt.Fprintf(w, text, a...)
	return nil
}

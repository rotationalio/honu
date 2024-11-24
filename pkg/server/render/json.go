package render

import (
	"encoding/json"
	"net/http"
)

const (
	jsonContentType = "application/json; charset=utf-8"
)

// JSON marshals the given interface object and writes it with the correct ContentType.
func JSON(code int, w http.ResponseWriter, obj any) error {
	w.Header().Set(ContentType, jsonContentType)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return err
	}
	return nil
}

package render

import (
	"net/http"
)

type Renderer interface {
	Render(code int, w http.ResponseWriter, obj any) error
}

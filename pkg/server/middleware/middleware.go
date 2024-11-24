package middleware

import "github.com/julienschmidt/httprouter"

type Middleware func(next httprouter.Handle) httprouter.Handle

func Chain(h httprouter.Handle, m ...Middleware) httprouter.Handle {
	if len(m) < 1 {
		return h
	}

	// Wrap the handler in the middleware in a reverse loop to preserve order
	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		if m[i] != nil {
			wrapped = m[i](wrapped)
		}
	}
	return wrapped
}

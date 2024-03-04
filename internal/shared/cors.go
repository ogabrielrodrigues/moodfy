package shared

import "net/http"

type CORS struct {
	origin string
}

func New(origin string) *CORS {
	return &CORS{origin}
}

func (c *CORS) Enable(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", c.origin)
}

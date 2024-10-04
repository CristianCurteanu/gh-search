package middlewares

import (
	"log"
	"net/http"
)

type Middleware interface {
	Execute(w http.ResponseWriter, r *http.Request) error
	GetFallback() http.HandlerFunc
}

type UseMiddleware struct {
	chain []Middleware
}

func (mw *UseMiddleware) Use(m Middleware) {
	mw.chain = append(mw.chain, m)
}

func (mw *UseMiddleware) Handle(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
	for _, fm := range mw.chain {
		err := fm.Execute(w, r)
		if fallback := fm.GetFallback(); err != nil && fallback != nil {
			fallback(w, r)
			return
		} else if fallback == nil {
			log.Printf("!!! WARNING: fallback is nil, falling through")
		}
	}

	h(w, r)
}

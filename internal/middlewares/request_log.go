package middlewares

import (
	"log"
	"net/http"
)

type RequestLog struct {
}

func NewRequestLog() Middleware {
	return &RequestLog{}
}

func (rl *RequestLog) Execute(w http.ResponseWriter, r *http.Request) error {
	log.Printf("[%s] %q", r.Method, r.URL.String())
	return nil
}

func (rl *RequestLog) GetFallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

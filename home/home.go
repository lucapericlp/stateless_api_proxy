package home

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const msg = "Hello"

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit Home")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func (h *Handlers) Files(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request for %s processed in %s\n", r.URL.Path, time.Now().Sub(startTime))
		next(w, r)
	}
}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	//mux.HandleFunc("/", h.Logger(h.Home))
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/static/", h.Logger(h.Files))
}

package home

import (
	"../keys"
	"../magictoken"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const msg = "Hello"

type TokenPOST struct {
	GithubToken string
	Scopes      []string
}

type TokenResponse struct {
	JWT string
}

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit Home")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var data TokenPOST
		err := decoder.Decode(&data)
		if err != nil {
			log.Fatal(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		ourKeys := keys.LoadKeys()
		ourJWT, err := magictoken.Create(data.GithubToken, data.Scopes, ourKeys)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error JWT", http.StatusBadRequest)
		}
		tokenResponse := &TokenResponse{
			JWT: ourJWT,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokenResponse)
		return
	}

	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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
	//mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/create", h.Create)
	mux.HandleFunc("/static/", h.Logger(h.Files))
}

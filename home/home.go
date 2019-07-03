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

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var data TokenPOST
		err := decoder.Decode(&data)
		if err != nil {
			h.logger.Fatal(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		ourKeys := keys.LoadKeys()
		ourJWT, err := magictoken.Create(data.GithubToken, data.Scopes, ourKeys)
		if err != nil {
			h.logger.Fatal(err)
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

func (h *Handlers) Verify(next http.HandlerFunc) http.HandlerFunc {
	//check validity of the token before passing to next step
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)
			var data TokenResponse
			err := decoder.Decode(&data)
			if err != nil {
				h.logger.Fatal(err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}

			ourKeys := keys.LoadKeys()
			ptToken, err := magictoken.Verify(data.JWT, ourKeys)
			if err != nil {
				http.Error(w, "Invalid proxy JWT supplied!", http.StatusBadRequest)
				h.logger.Println(err)
				return
			}
			//fmt.Println(ptToken, time.Now())
			h.logger.Println(ptToken, time.Now())
			next(w, r)
		}
		http.Error(w, "POST Route", http.StatusBadRequest)
		return
	}
}

func (h *Handlers) Api(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit API")
}

func (h *Handlers) Files(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

//func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		startTime := time.Now()
//		defer h.logger.Printf("request for %s processed in %s\n", r.URL.Path, time.Now().Sub(startTime))
//		next(w, r)
//	}
//}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	//mux.HandleFunc("/", h.Logger(h.Home))
	//mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/create", h.Create)
	mux.HandleFunc("/api", h.Verify(h.Api))
	mux.HandleFunc("/static/", h.Files) //h.Logger(h.Files))
}

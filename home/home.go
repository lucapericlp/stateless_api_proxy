package home

import (
	"../keys"
	"../magictoken"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type OriginalToken struct {
	GithubToken string
	Scopes      []string
}

type TokenResponse struct {
	JWT string
}

type Handlers struct {
	logger *log.Logger
}

//Create JWT from provided GithubToken & Scopes along with issuedat and expiresat timestamps which then consume env priv and pub keys.
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var data OriginalToken
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
		//what are they trying to do?
		h.logger.Println(r.Method, r.URL.Path)
		method, path := r.Method, r.URL.Path

		jwt := r.Header.Get("Authorization")
		if strings.HasPrefix(jwt, "Bearer ") {
			jwt = strings.TrimPrefix(jwt, "Bearer ")
		} else {
			http.Error(w, "JWT not supplied via Authorization header", http.StatusBadRequest)
			h.logger.Println("JWT not in Auth header")
			return
		}

		ourKeys := keys.LoadKeys()
		proxyToken, err := magictoken.Verify(jwt, ourKeys)
		if err != nil {
			http.Error(w, "Invalid proxy JWT supplied!", http.StatusBadRequest)
			h.logger.Println(err)
			return
		}

		authorisedScopes := strings.Join(proxyToken.Scopes, ", ")
		if !proxyToken.ValidateRequest(method, path) {

			var errMsg strings.Builder
			errMsg.WriteString("Unauthorised scope according to Github proxy. Authorised scopes: ")
			errMsg.WriteString(authorisedScopes)

			http.Error(w, errMsg.String(), http.StatusUnauthorized)
			h.logger.Println(errMsg.String())
			return
		}

		h.logger.Printf("Verified %v req for %v with scopes %v", method, path, authorisedScopes)
		next(w, r)
		return
	}
}

func (h *Handlers) Api(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit API")
}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create", h.Create)
	mux.HandleFunc("/api/", h.Verify(h.Api))
}

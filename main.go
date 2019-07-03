package main

import (
	"./home"
	"./keys"
	"./magictoken"
	"./server"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	CertFile    = os.Getenv("TEST_CERT_FILE")
	KeyFile     = os.Getenv("TEST_KEY_FILE")
	ServiceAddr = os.Getenv("TEST_SERVER_ADDR")
)

func main() {
	logger := log.New(os.Stdout, "test", log.LstdFlags|log.Lshortfile)
	h := home.NewHandlers(logger)

	//ourKeys := keys.LoadKeys()
	//proxyToken, _ := magictoken.Create("abc123", []string{"GET /user", "GET /repos"}, ourKeys)

	//ptToken, _ := magictoken.Verify(proxyToken, ourKeys)
	//fmt.Println(ptToken)

	mux := http.NewServeMux()
	h.SetupRoutes(mux)
	srv := server.New(mux, ServiceAddr)

	logger.Println("server starting")
	err := srv.ListenAndServeTLS(CertFile, KeyFile)
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

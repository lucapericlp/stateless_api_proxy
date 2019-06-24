package main

import (
	"./home"
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

	//	magictoken.Encrypt()
	token := magictoken.Create("abc123", [2]string{"1", "2"})
	fmt.Println(token)

	mux := http.NewServeMux()
	h.SetupRoutes(mux)
	srv := server.New(mux, ServiceAddr)

	logger.Println("server starting")
	err := srv.ListenAndServeTLS(CertFile, KeyFile)
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

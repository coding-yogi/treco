// Package server
package server

import (
	"fmt"
	"log"
	"net/http"
	"treco/storage"
)

// Starts the server mode
func Start(port int) {
	var err error

	// Connect to storage
	err = storage.New()
	if err != nil {
		log.Fatal(err)
	}

	handler := storage.Handler()
	defer func() {
		_ = (*handler).Close()
	}()

	//DB setup
	//err = (*handler).Setup(dbEntities...)
	//exitOnError(err)

	// Define http handler
	var publisherHandler PublishHandler
	http.HandleFunc("/v1/publish/report", publisherHandler.ServeHTTP)

	// start server
	log.Printf("Starting server on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

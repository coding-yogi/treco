// Package server
package server

import (
	"fmt"
	"log"
	"net/http"
	"treco/conf"
	"treco/model"
	"treco/storage"
)

var DBEntities = []interface{}{&model.SuiteResult{}, &model.ScenarioResult{}, &model.Scenario{}, &model.Feature{}}

// Starts the server mode
func Start(cfgFile string, port int) {
	var err error

	// check config file
	if cfgFile != "" {
		if err := conf.LoadEnvFromFile(cfgFile); err != nil {
			log.Fatalf("error occured while loading from config %v\n", err)
			return
		}
	} else {
		log.Println("no config file path set")
	}

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
	err = (*handler).Setup(DBEntities...)
	if err != nil {
		log.Fatal(err)
	}

	// Define http handler
	var publisherHandler PublishHandler
	http.HandleFunc("/v1/publish/report", publisherHandler.ServeHTTP)

	// start server
	log.Printf("Starting server on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

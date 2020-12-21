package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

var requiredParams= []string{ReportFormat, TestType, BuildID, Service}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs as a web server",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/treco/v1/publish", publishHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

type Error struct {
	Code        int
	Description string
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	// Validate Method
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	// Validate Headers
	if err := validateHeaders(r.Header); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse parameters
	if err := r.ParseMultipartForm(1000000); err != nil {
		sendErrorResponse(w, "unable to parse form data", http.StatusBadRequest)
		return
	}

	// Validate parameters
	if err := validateFormParameters(r); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg = config{
		reportFile: "",
		reportFormat: r.FormValue(strings.ToLower(ReportFormat)),
		testType: r.FormValue(strings.ToLower(TestType)),
		build: r.FormValue(strings.ToLower(BuildID)),
		service: r.FormValue(strings.ToLower(Service)),
	}

	// Validate param values
	if err:= validateParams(cfg); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read file from report_file
	reportFile, _, err := r.FormFile(strings.ToLower(ReportFile))
	if err != nil {
		sendErrorResponse(w, "unable to retrieve report file", http.StatusBadRequest)
		return
	}

	defer reportFile.Close()

	// Process file
	if err := process(cfg, reportFile); err != nil {
		log.Println("error while processing file: " + err.Error())
		sendErrorResponse(w, "unable to process the request", http.StatusInternalServerError)
		return
	}

	log.Println("results uploaded successfully")
	w.WriteHeader(http.StatusOK)
	return
}

func sendErrorResponse(w http.ResponseWriter,  message string, code int) {
	b,_ := json.Marshal(Error{
		Code: code,
		Description: message,
	})

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(b)
}

func validateHeaders(h http.Header) error {
	//check content-type
	expectedContentType := "multipart/form-data"
	if !strings.Contains(h.Get("content-type"), "multipart/form-data") {
		return fmt.Errorf("invalid content-type, expected: %v", expectedContentType)
	}

	return nil
}

func validateFormParameters(r *http.Request) error {
	missingParams := make([]string, 0)
	for _, param := range requiredParams {
		param := strings.ToLower(param)
		if r.FormValue(param) == "" {
			missingParams = append(missingParams, param)
		}
	}

	if len(missingParams) > 0 {
		return fmt.Errorf("missing params: %v", strings.Join(missingParams, ", "))
	}
	return nil
}










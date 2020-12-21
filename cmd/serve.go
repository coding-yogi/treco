package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

type Error struct {
	Code        int
	Description string
}

var port int
var requiredParams= []string{ReportFormat, TestType, BuildID, Service}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs as a web server",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/treco/v1/publish/report", publishHandler)
		log.Printf("Starting server on port %v\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
	},
}

func init() {
	flags := serveCmd.Flags()
	flags.IntVarP(&port, "port", "p", 8080, "port for server to run")
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if status, err := validatePublishRequest(r); err != nil {
		sendErrorResponse(w, err.Error(), status)
		return
	}

	// Read file from report_file
	reportFile, _, err := r.FormFile(strings.ToLower(ReportFile))
	if err != nil {
		sendErrorResponse(w, "unable to retrieve report file", http.StatusBadRequest)
		return
	}

	defer reportFile.Close()

	cfg = config{
		reportFile: "",
		reportFormat: r.FormValue(strings.ToLower(ReportFormat)),
		testType: r.FormValue(strings.ToLower(TestType)),
		build: r.FormValue(strings.ToLower(BuildID)),
		service: r.FormValue(strings.ToLower(Service)),
	}

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

func validatePublishRequest(r *http.Request) (int, error) {
	// Validate Method
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, fmt.Errorf("")
	}

	// Validate content-type
	expectedContentType := "multipart/form-data"
	if !strings.Contains(r.Header.Get("content-type"), expectedContentType) {
		return http.StatusBadRequest, fmt.Errorf("invalid content-type, expected: %v", expectedContentType)
	}

	// Validate parameters
	missingParams := make([]string, 0)
	for _, param := range requiredParams {
		param := strings.ToLower(param)
		if r.FormValue(param) == "" {
			missingParams = append(missingParams, param)
		}
	}

	if len(missingParams) > 0 {
		return http.StatusBadRequest, fmt.Errorf("missing params: %v", strings.Join(missingParams, ", "))
	}

	// Validate param values
	if err:= validateParams(r.FormValue(strings.ToLower(TestType)),
		r.FormValue(strings.ToLower(ReportFormat))); err != nil {
		return http.StatusBadRequest, err
	}

	return 0, nil
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










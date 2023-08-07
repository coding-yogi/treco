package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"treco/conf"
	"treco/model"
	"treco/report"
	"treco/storage"
)

const (
	BuildID      = "CI_JOB_ID"
	Environment  = "ENVIRONMENT"
	Jira         = "JIRA_PROJECT"
	ReportFile   = "REPORT_FILE"
	ReportFormat = "REPORT_FORMAT"
	Service      = "SERVICE_NAME"
	TestType     = "TEST_TYPE"
	Coverage     = "COVERAGE"

	expectedContentType = "multipart/form-data"
)

var (
	validTestTypes     = [...]string{"unit", "contract", "integration", "e2e"}
	validReportFormats = [...]string{"junit"}

	errInvalidTestType       = "test type %v is invalid, should be one of %v"
	errInvalidReportFormats  = "report format %v is invalid, should be one of %v"
	errCoverageValueNotFloat = "coverage value should be a floating number"

	// RequiredParams ...
	RequiredParams = [...]string{BuildID, Environment, Jira, ReportFormat, Service, TestType}
)

// Error struct having code and description
type Error struct {
	Code        int
	Description string
}

// PublishHandler ...
type PublishHandler struct {
}

// ServerHTTP ...
func (p PublishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if status, err := validatePublishRequest(r); err != nil {
		sendErrorResponse(w, err, err.Error(), status)
		return
	}

	rf, err := readFileFromRequest(r)
	if err != nil {
		sendErrorResponse(w, err, "unable to retrieve report file", http.StatusBadRequest)
		return
	}

	cfg := conf.Config{
		Build:        r.FormValue(strings.ToLower(BuildID)),
		Environment:  r.FormValue(strings.ToLower(Environment)),
		Jira:         r.FormValue(strings.ToLower(Jira)),
		Service:      r.FormValue(strings.ToLower(Service)),
		ReportFormat: r.FormValue(strings.ToLower(ReportFormat)),
		TestType:     r.FormValue(strings.ToLower(TestType)),
		Coverage:     r.FormValue(strings.ToLower(Coverage)),
	}

	// Process file
	if err := Process(cfg, rf); err != nil {
		log.Println("error processing: " + err.Error())
		sendErrorResponse(w, err, "unable to process the request", http.StatusInternalServerError)
		return
	}

	log.Println("results uploaded successfully")
	w.WriteHeader(http.StatusOK)
}

// Read the report file from request
func readFileFromRequest(r *http.Request) (io.Reader, error) {
	reportFile, _, err := r.FormFile(strings.ToLower(ReportFile))
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = reportFile.Close()
	}()

	return reportFile, nil
}

// Validate incoming publish request
func validatePublishRequest(r *http.Request) (int, error) {
	// Validate Method
	if r.Method != "POST" {
		return http.StatusMethodNotAllowed, fmt.Errorf("")
	}

	// Validate content-type
	if !strings.Contains(r.Header.Get("content-type"), expectedContentType) {
		return http.StatusBadRequest, fmt.Errorf("invalid content-type, expected: %s", expectedContentType)
	}

	// Validate parameters
	missingParams := make([]string, 0, len(RequiredParams))
	for _, param := range RequiredParams {
		lparam := strings.ToLower(param)
		if r.FormValue(lparam) == "" {
			missingParams = append(missingParams, lparam)
		}
	}

	if len(missingParams) > 0 {
		return http.StatusBadRequest, fmt.Errorf("missing params: %v", strings.Join(missingParams, ", "))
	}

	// Validate param values
	testType := r.FormValue(strings.ToLower(TestType))
	reportFormat := r.FormValue(strings.ToLower(ReportFormat))
	coverage := r.FormValue(strings.ToLower(Coverage))
	if err := ValidateParams(testType, reportFormat, coverage); err != nil {
		return http.StatusBadRequest, err
	}

	return 0, nil
}

// send error response
func sendErrorResponse(w http.ResponseWriter, err error, description string, code int) {
	log.Println(err)

	b, _ := json.Marshal(Error{
		Code:        code,
		Description: description,
	})

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(b)
}

// proces the request
func Process(cfg conf.Config, f io.Reader) error {
	var err error
	coverage, _ := strconv.ParseFloat(cfg.Coverage, 64)

	// Create base result
	data := &model.Data{
		Jira:         cfg.Jira,
		ReportFormat: cfg.ReportFormat,
		SuiteResult: model.SuiteResult{
			Build:       cfg.Build,
			Environment: cfg.Environment,
			Service:     strings.ToLower(cfg.Service),
			TestType:    strings.ToLower(cfg.TestType),
			Coverage:    coverage,
		},
	}

	// Transform file data into required format
	err = report.Parse(f, data)
	if err != nil {
		return err
	}

	// Write to storage
	dbh := storage.Handler()
	err = data.Save(dbh)
	if err != nil {
		return err
	}

	return nil
}

// Validates if Test Type, Report Type and Coverage values are valid
func ValidateParams(testType, reportType, coverage string) error {
	//check for valid test type
	if !isValid(testType, validTestTypes[:]) {
		return fmt.Errorf(errInvalidTestType, testType, validTestTypes)
	}

	//check for valid test report format
	if !isValid(reportType, validReportFormats[:]) {
		return fmt.Errorf(errInvalidReportFormats, reportType, validReportFormats)
	}

	//check coverage is in float
	if _, err := strconv.ParseFloat(coverage, 64); err != nil {
		return fmt.Errorf(errCoverageValueNotFloat)
	}

	return nil
}

// Checks that value exist in an array of valid values
func isValid(value string, values []string) bool {
	value = strings.ToLower(value)
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}

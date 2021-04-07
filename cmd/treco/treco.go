/*
Package cmd runs tool as command line util
*/
package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"treco/model"
	"treco/report"
	"treco/storage"

	"github.com/spf13/cobra"
)

// Command Line params
const (
	BuildID      = "CI_JOB_ID"
	Environment  = "ENVIRONMENT"
	Jira         = "JIRA_PROJECT"
	ReportFile   = "REPORT_FILE"
	ReportFormat = "REPORT_FORMAT"
	Service      = "SERVICE_NAME"
	TestType     = "TEST_TYPE"
)

var cfg config

type config struct {
	Build        string
	Environment  string
	Jira         string
	ReportFile   string
	ReportFormat string
	Service      string
	TestType     string
}

var rootCmd = &cobra.Command{
	Use:   "treco",
	Short: "Test Report Collector",
}

// Execute root commans
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(serveCmd)
}

var (
	validTestTypes     = [...]string{"unit", "contract", "integration", "e2e"}
	validReportFormats = [...]string{"junit"}

	errInvalidTestType      = fmt.Errorf("test type should be one of %v", validTestTypes)
	errInvalidReportFormats = fmt.Errorf("report format should be one of %v", validReportFormats)
)

func validateParams(testType, reportType string) error {
	//check for valid test type
	if !isValid(testType, validTestTypes[:]) {
		return errInvalidTestType
	}

	//check for valid test report format
	if !isValid(reportType, validReportFormats[:]) {
		return errInvalidReportFormats
	}

	return nil
}

func isValid(value string, values []string) bool {
	value = strings.ToLower(value)
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}

func exitOnError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func process(cfg config, f io.Reader) error {
	var err error

	// Create base result
	data := &model.Data{
		Jira:         cfg.Jira,
		ReportFormat: cfg.ReportFormat,
		SuiteResult: model.SuiteResult{
			Build:       cfg.Build,
			Environment: cfg.Environment,
			Service:     strings.ToLower(cfg.Service),
			TestType:    strings.ToLower(cfg.TestType),
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

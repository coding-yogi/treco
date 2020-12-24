package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"strings"
	"treco/model"
	"treco/report"
)

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
	build        string
	environment  string
	jira         string
	reportFile   string
	reportFormat string
	service      string
	testType     string
}

var validTestTypes = []string{"unit", "contract", "integration", "e2e"}
var validReportFormats = []string{"junit"}

var (
	ErrInvalidTestType     = fmt.Errorf("test type should be one of %v", validTestTypes)
	ErrInvalidReportFormat = fmt.Errorf("report format should be one of %v", validReportFormats)
)

var rootCmd = &cobra.Command{
	Use:   "treco",
	Short: "Test Report Collector",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(serveCmd)
}

func validateParams(testType, reportType string) error {
	//check for valid test type
	if !isValid(testType, validTestTypes) {
		return ErrInvalidTestType
	}

	//check for valid test report format
	if !isValid(reportType, validReportFormats) {
		return ErrInvalidReportFormat
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

func process(cfg *config, f *io.Reader) error {
	var err error

	// Create base result
	data := model.Data{
		Jira:         cfg.jira,
		ReportFormat: cfg.reportFormat,
		SuiteResult: model.SuiteResult{
			Build:       cfg.build,
			Environment: cfg.environment,
			Service:     strings.ToLower(cfg.service),
			TestType:    strings.ToLower(cfg.testType),
		},
	}

	// Transform file data into required format
	err = report.Parse(f, &data)
	if err != nil {
		return err
	}

	// Write to storage
	err = data.Save()
	if err != nil {
		return err
	}

	return nil
}

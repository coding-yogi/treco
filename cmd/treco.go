package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"strings"
	"time"
	"treco/model"
	"treco/report"
	"treco/storage"
)

const (
	ReportFile   = "REPORT_FILE"
	ReportFormat = "REPORT_FORMAT"
	Service      = "SERVICE_NAME"
	TestType     = "TEST_TYPE"
	BuildID      = "CI_JOB_ID"
	Jira         = "JIRA_PROJECT"
)

var cfg config

type config struct {
	reportFile   string
	reportFormat string
	service      string
	testType     string
	build        string
	jira         string
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

func createDBSchema(dbh storage.DBHandler) error {
	switch db := dbh.(type) {
	case storage.Postgres:
		db.CreateSchema([]interface{}{
			(*model.SuiteResult)(nil),
			(*model.Scenario)(nil),
			(*model.Feature)(nil),
			(*model.ScenarioResult)(nil),
		})
	}

	return nil
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

func process(cfg config, f io.Reader) error {

	var err error

	// Create base result
	data := model.Data{
		Jira:         cfg.jira,
		ReportFormat: cfg.reportFormat,
		SuiteResult: model.SuiteResult{
			Build:      cfg.build,
			Service:    strings.ToLower(cfg.service),
			TestType:   strings.ToLower(cfg.testType),
			ExecutedAt: time.Now(),
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

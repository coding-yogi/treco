package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"time"
	"treco/report"
	"treco/model"
)

const (
	REPORT_FILE   = "REPORT_FILE"
	REPORT_FORMAT = "REPORT_FORMAT"
	SERVICE       = "SERVICE_NAME"
	TEST_TYPE     = "TEST_TYPE"
	BUILD         = "CI_JOB_ID"
)

var config Config

type Config struct {
	ReportFile   string
	ReportFormat string
	Service      string
	TestType     string
	Build        string
}

var validTestTypes = []string{"unit", "contract", "integration", "e2e"}
var validReportFormats =[]string{"junit"}

var rootCmd = &cobra.Command{
	Use:   "treco",
	Short: "Commandline tool to push your test stats to datasource",
	Run: func(cmd *cobra.Command, args []string) {
		//validate flags
		if err := validateFlags(config); err!=nil {
			return
		}

		//check for report file
		reportFile, err := open(config.ReportFile)
		if err != nil {
			return
		}

		//transform file data into required format
		transform(reportFile, config.ReportFormat)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	flags := rootCmd.Flags()
	flags.StringVarP(&config.ReportFile, "report", "r", os.Getenv(REPORT_FILE), "input file containing test reports")
	flags.StringVarP(&config.ReportFormat, "format", "f", os.Getenv(REPORT_FORMAT), "report of report file")
	flags.StringVarP(&config.Service, "service", "s", os.Getenv(SERVICE), "service name")
	flags.StringVarP(&config.TestType, "type", "t", os.Getenv(TEST_TYPE), "type of tests executed. 'unit', 'contract', 'integration' or 'e2e")
	flags.StringVarP(&config.Build, "build", "b", os.Getenv(BUILD), "CI build name or number to uniquely identify the build")
}

func validateFlags(config Config) error {
	//check for empty flags
	if config.ReportFile == "" || config.ReportFormat == "" || config.Service == "" || config.TestType == "" || config.Build == "" {
		return fmt.Errorf("\nmissing arguments, please run `treco --help` for more info\n" +
			"\nyou can also supply arguments via following ENVIRONMENT variables\n" +
			"export %s=\n" +
			"export %s=\n" +
			"export %s=\n" +
			"export %s=\n" +
			"export %s=\n", REPORT_FILE, REPORT_FORMAT, SERVICE, TEST_TYPE, BUILD)
	}

	//check for valid test type
	if !isValid(config.TestType, validTestTypes) {
		return fmt.Errorf("test type should be one of %v. %v is not a valid test type\n", validTestTypes, config.TestType )
	}

	//check for valid test report format
	if !isValid(config.ReportFormat, validReportFormats) {
		return fmt.Errorf("report should be one of %v. %v is not a valid report report\n", validReportFormats, config.ReportFormat )
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

func open(f string) (*os.File, error) {
	file, err := os.OpenFile(f, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	return file, err
}

func transform(r io.Reader, rf string) error {
	var parser report.ReportParser
	var err error

	result := model.Result{
		Build:    config.Build,
		Service:  config.Service,
		TestType: config.TestType,
		Time:     time.Now().Unix(),
	}

	switch strings.ToLower(rf) {
	case "junit":
		parser = report.JunitXmlParser{}
		err = parser.Parse(r, &result)
	default:
		err = errors.New("invalid report " + rf)
	}

	return err
}
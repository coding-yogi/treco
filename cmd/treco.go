package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
	"treco/model"
	"treco/report"
	"treco/storage"
)

const (
	ReportFile   = "ReportFile"
	ReportFormat = "ReportFormat"
	Service      = "SERVICE_NAME"
	TestType     = "TestType"
	BuildID      = "CI_JOB_ID"
)

var cfg config

type config struct {
	reportFile   string
	reportFormat string
	service      string
	testType     string
	build        string
}

var validTestTypes = []string{"unit", "contract", "integration", "e2e"}
var validReportFormats = []string{"junit"}

var rootCmd = &cobra.Command{
	Use:   "treco",
	Short: "Commandline tool to push your test stats to datasource",
	Run: func(cmd *cobra.Command, args []string) {
		//validate flags
		err := validateFlags(cfg)
		exitOnError(err)

		//check for report file
		reportFile, err := os.OpenFile(cfg.reportFile, os.O_RDONLY, 0644)
		exitOnError(err)
		defer reportFile.Close()

		//connect to storage
		executor, err := storage.New()
		exitOnError(err)
		defer executor.Close()

		result := model.Result{
			DbHandler:  executor,
			Build:      cfg.build,
			Service:    strings.ToLower(cfg.service),
			TestType:   strings.ToLower(cfg.testType),
			ExecutedAt: time.Now(),
		}

		//transform file data into required format
		err = report.Parse(reportFile, cfg.reportFormat, &result)
		exitOnError(err)

		//write to storage
		err = result.Save()
		exitOnError(err)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	flags := rootCmd.Flags()
	flags.StringVarP(&cfg.reportFile, "report", "r", os.Getenv(ReportFile), "input file containing test reports")
	flags.StringVarP(&cfg.reportFormat, "format", "f", os.Getenv(ReportFormat), "report of report file")
	flags.StringVarP(&cfg.service, "service", "s", os.Getenv(Service), "service name")
	flags.StringVarP(&cfg.testType, "type", "t", os.Getenv(TestType), "type of tests executed. 'unit', 'contract', 'integration' or 'e2e")
	flags.StringVarP(&cfg.build, "build", "b", os.Getenv(BuildID), "CI build name or number to uniquely identify the build")
}

func validateFlags(cfg config) error {
	//check for empty flags
	log.Println("validating parameters")
	if cfg.reportFile == "" || cfg.reportFormat == "" || cfg.service == "" || cfg.testType == "" || cfg.build == "" {
		return fmt.Errorf("\nmissing arguments, please run `treco --help` for more info\n"+
			"\nyou can also supply arguments via following ENVIRONMENT variables\n"+
			"export %s=\n"+
			"export %s=\n"+
			"export %s=\n"+
			"export %s=\n"+
			"export %s=\n", ReportFile, ReportFormat, Service, TestType, BuildID)
	}

	//check for valid test type
	if !isValid(cfg.testType, validTestTypes) {
		return fmt.Errorf("test type should be one of %v. %v is not a valid test type\n", validTestTypes, cfg.testType)
	}

	//check for valid test report format
	if !isValid(cfg.reportFormat, validReportFormats) {
		return fmt.Errorf("report should be one of %v. %v is not a valid report report\n", validReportFormats, cfg.reportFormat)
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

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"treco/storage"

	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Runs as a command line tool",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Connect to storage
		err = storage.New()
		exitOnError(err)

		handler := storage.Handler()
		defer func() {
			_ = (*handler).Close()
		}()

		//validate flags
		err = validateFlags(cfg)
		exitOnError(err)

		//check for report file
		reportFile, err := os.OpenFile(cfg.ReportFile, os.O_RDONLY, 0600)
		exitOnError(err)
		defer func() {
			_ = reportFile.Close()
		}()

		// Process file
		var rf io.Reader = reportFile
		err = process(cfg, rf)
		exitOnError(err)

		log.Println("results uploaded successfully")
	},
}

func init() {
	flags := collectCmd.Flags()

	flags.StringVarP(&cfg.Build, "Build", "b", os.Getenv(BuildID), "CI Build name or number to uniquely identify the Build")
	flags.StringVarP(&cfg.Environment, "Environment", "e", os.Getenv(Environment), "Environment on which the Build is executed")
	flags.StringVarP(&cfg.Jira, "Jira", "j", os.Getenv(Jira), "Jira project name")
	flags.StringVarP(&cfg.ReportFile, "report", "r", os.Getenv(ReportFile), "input file containing test reports")
	flags.StringVarP(&cfg.ReportFormat, "format", "f", os.Getenv(ReportFormat), "report of report file")
	flags.StringVarP(&cfg.Service, "Service", "s", os.Getenv(Service), "Service name")
	flags.StringVarP(&cfg.TestType, "type", "t", os.Getenv(TestType), "type of tests executed. 'unit', 'contract', 'integration' or 'e2e")
}

var (
	errMissingArguments = fmt.Errorf("\nmissing arguments, please run `treco --help` for more info\n"+
		"\nyou can also supply arguments via following ENVIRONMENT variables\n"+
		"%s, %s, %s, %s, %s, %s, %s ", BuildID, Environment, Jira, ReportFile, ReportFormat, Service, TestType)
)

func validateFlags(cfg config) error {
	//check for empty flags
	log.Println("validating parameters")
	if cfg.ReportFile == "" || cfg.ReportFormat == "" || cfg.Service == "" || cfg.TestType == "" || cfg.Build == "" || cfg.Jira == "" || cfg.Environment == "" {
		return errMissingArguments
	}

	return validateParams(cfg.TestType, cfg.ReportFormat)
}

package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"treco/conf"
	"treco/server"
	"treco/storage"

	"github.com/spf13/cobra"
)

func newCollectCommand() *cobra.Command {
	var cfg conf.Config

	collectCmd := &cobra.Command{
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

			//DB setup
			err = (*handler).Setup(server.DBEntities...)
			exitOnError(err)

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
			err = server.Process(cfg, rf)
			exitOnError(err)

			log.Println("results uploaded successfully")
		},
	}

	flags := collectCmd.Flags()

	flags.StringVarP(&cfg.Build, "build", "b", os.Getenv(server.BuildID), "CI Build name or number to uniquely identify the Build")
	flags.StringVarP(&cfg.Environment, "environment", "e", os.Getenv(server.Environment), "Environment on which the Build is executed")
	flags.StringVarP(&cfg.Jira, "jira", "j", os.Getenv(server.Jira), "Jira project name")
	flags.StringVarP(&cfg.ReportFile, "report", "r", os.Getenv(server.ReportFile), "input file containing test reports")
	flags.StringVarP(&cfg.ReportFormat, "format", "f", os.Getenv(server.ReportFormat), "report of report file")
	flags.StringVarP(&cfg.Service, "service", "s", os.Getenv(server.Service), "Service name")
	flags.StringVarP(&cfg.TestType, "type", "t", os.Getenv(server.TestType), "type of tests executed. 'unit', 'contract', 'integration' or 'e2e")
	flags.StringVarP(&cfg.Coverage, "coverage", "c", os.Getenv(server.Coverage), "statement level code coverage")

	return collectCmd
}

var (
	errMissingArguments = fmt.Errorf("\nmissing arguments, please run `treco --help` for more info\n"+
		"\nyou can also supply arguments via following ENVIRONMENT variables\n"+
		"%v ", server.RequiredParams)
)

// validate flags sent to collect command
func validateFlags(cfg conf.Config) error {
	//check for empty flags
	log.Println("validating parameters")
	if cfg.ReportFile == "" || cfg.ReportFormat == "" || cfg.Service == "" || cfg.TestType == "" || cfg.Build == "" ||
		cfg.Jira == "" || cfg.Environment == "" || cfg.Coverage == "" {
		return errMissingArguments
	}

	return server.ValidateParams(cfg.TestType, cfg.ReportFormat, cfg.Coverage)
}

// exits ith fatal error
func exitOnError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

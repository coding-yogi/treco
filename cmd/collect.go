package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Runs as a command line tool",
	Run: func(cmd *cobra.Command, args []string) {
		//validate flags
		err := validateFlags(cfg)
		exitOnError(err)

		//check for report file
		reportFile, err := os.OpenFile(cfg.reportFile, os.O_RDONLY, 0644)
		exitOnError(err)
		defer reportFile.Close()

		// Process file
		err = process(cfg, reportFile)
		exitOnError(err)

		log.Println("results uploaded successfully")
	},
}

func init() {
	flags := collectCmd.Flags()
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

	return validateParams(cfg.testType, cfg.reportFormat)
}

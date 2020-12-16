package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"treco/format"
	"treco/model"
	"time"
)

var (
	input        string
	reportFormat string
	service      string
	testType     string
	build        string
)

const emptyString = ""

var rootCmd = &cobra.Command{
	Use:   "treco",
	Short: "Commandline tool to push your test stats to datasource",
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateEmptyFlags(); err != nil {
			fmt.Println(err.Error())
			return
		}

		if err := validateReportFormat(reportFormat); err != nil {
			fmt.Println(err.Error())
			return
		}

		if err := validateTestType(testType); err != nil {
			fmt.Println(err.Error())
			return
		}

		reportFile, err := open(input)
		if err != nil {
			return
		}

		transform(reportFile, reportFormat)

	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	flags := rootCmd.Flags()
	flags.StringVarP(&input, "input", "i", emptyString, "input file containing test reports")
	flags.StringVarP(&reportFormat, "format", "f", emptyString, "format of report file")
	flags.StringVarP(&service, "service", "s", emptyString, "service name")
	flags.StringVarP(&testType, "type", "t", emptyString, "type of tests executed. 'unit', 'contract', 'integration' or 'e2e")
	flags.StringVarP(&build, "build", "b", emptyString, "build name or number from CI tool")
}

func validateEmptyFlags() error {
	input = strings.TrimSpace(input)
	reportFormat = strings.TrimSpace(reportFormat)
	service = strings.TrimSpace(service)
	testType = strings.TrimSpace(testType)
	build = strings.TrimSpace(build)

	if input == "" || reportFormat == "" || service == "" || testType == "" || build == "" {
		return errors.New("not all parameters supplied, please run `treco --help` for more info")
	}

	return nil
}

func validateReportFormat(format string) error {
	allowedFormats := []string{"junit"}
	format = strings.ToLower(format)

	for _, f := range allowedFormats {
		if f == format {
			return nil
		}
	}
	return fmt.Errorf("report format should be either of %v", allowedFormats)
}

func validateTestType(testType string) error {
	allowedTestTypes := []string{"unit", "contract", "integration", "e2e"}
	testType = strings.ToLower(testType)

	for _, tt := range allowedTestTypes {
		if tt == testType {
			return nil
		}
	}
	return fmt.Errorf("type of tests should be either of %v", allowedTestTypes)
}

func open(f string) (*os.File, error) {
	file, err := os.OpenFile(f, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("error reading file " + f + err.Error())
	}

	return file, err
}

func transform(r io.Reader, rf string) error {
	var parser format.Parser
	var err error

	result := model.Result{
		Build: build,
		Service: service,
		TestType: testType,
		Time: time.Now().Unix(),
	}

	switch strings.ToLower(rf) {
	case "junit":
		parser = format.JunitXmlParser{}
		err = parser.Parse(r, &result)
	default:
		err = errors.New("invalid format " + rf)
	}

	return err
}
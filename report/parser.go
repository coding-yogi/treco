/*
Package report handles parsing of different report types
*/
package report

import (
	"fmt"
	"io"
	"strings"
	"treco/model"
)

var (
	errInvalidReportType = "invalid report type: %v"
)

type Parser interface {
	parse(r io.Reader, result *model.Data) error
}

// Parse parses data from provided reader
func Parse(r io.Reader, data *model.Data) error {
	var parser Parser
	var err error

	rf := strings.ToLower(data.ReportFormat)

	switch rf {
	case "junit":
		parser = junitXMLParser{}
		err = parser.parse(r, data)
	default:
		err = fmt.Errorf(errInvalidReportType, rf)
	}

	return err
}

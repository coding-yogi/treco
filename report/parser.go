package report

import (
	"fmt"
	"io"
	"strings"
	"treco/model"
)

type parser interface {
	parse(r io.Reader, result *model.Data) error
}

func Parse(r io.Reader, data *model.Data) error {
	var parser parser
	var err error

	rf := strings.ToLower(data.ReportFormat)

	switch rf {
	case "junit":
		parser = junitXmlParser{}
		err = parser.parse(r, data)
	default:
		err = fmt.Errorf("invalid report type: %v " , rf)
	}

	return err
}

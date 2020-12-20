package report

import (
	"errors"
	"io"
	"strings"
	"treco/model"
)

type parser interface {
	parse(r io.Reader, result *model.Result) error
}

func Parse(r io.Reader, rf string, result *model.Result) error {
	var parser parser
	var err error

	switch strings.ToLower(rf) {
	case "junit":
		parser = junitXmlParser{}
		err = parser.parse(r, result)
	default:
		err = errors.New("invalid report type " + rf)
	}

	return err
}

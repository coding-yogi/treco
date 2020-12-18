package report

import (
	"io"
	"treco/model"
)

type ReportParser interface {
	Parse(r io.Reader, result *model.Result) error
}

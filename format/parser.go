package format

import (
	"io"
	"treco/model"
)

type Parser interface {
	Parse(r io.Reader, result *model.Result) error
}

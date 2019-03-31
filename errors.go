package tools

import (
	"bytes"
	"github.com/pkg/errors"
)

func ErrorsDump(err error) string {
	return errors.Cause(err).Error()
}

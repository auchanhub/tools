package tools

import (
	"github.com/pkg/errors"
)

func ErrorsDump(err error) string {
	return errors.Cause(err).Error()
}

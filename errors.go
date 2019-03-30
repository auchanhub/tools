package tools

import (
	"bytes"
	"github.com/pkg/errors"
)

func ErrorsDump(err error) string {
	str := bytes.NewBuffer([]byte{})

	errors.Wrap(str, err)

	return str.String()
}

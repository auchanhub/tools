package fhttp

import (
	"net/http"
	"os"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

func RequestParseParams(req *http.Request, params interface{}) (err error) {
	if req == nil {
		err = errors.Wrap(os.ErrInvalid, "request is nil")
		return
	}

	if err = req.ParseForm(); err != nil {
		err = errors.Wrap(os.ErrInvalid, "failed to parse request form params")
		return
	}

	var requestParams map[string][]string

	switch req.Method {
	case "GET":
		requestParams = req.URL.Query()

	case "POST":
		requestParams = req.Form
	}

	if err = schema.NewDecoder().Decode(params, requestParams); err != nil {
		err = errors.Wrap(os.ErrInvalid, "failed to decode request params")
		return
	}

	return
}

func BindHandlerParams(req *http.Request, data interface{}) interface{} {
	if err := RequestParseParams(req, data); err != nil {
		//	TODO: log errors
	}

	return data
}
package fhttp

import (
	"net/http"
	"encoding/json"
    "github.com/pkg/errors"
)

// The http handler is with result data which should convert to JSON
type JsonFuncHandler func(w http.ResponseWriter, req *http.Request) (interface{}, error)

func JsonFuncHandle(handler JsonFuncHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")

		d, err := handler(w, req)
		if err != nil {
			// TODO: handle error - should write to log or other reaction
			return
		}

		response, err := json.Marshal(d)
		if err != nil {
			// TODO: handle error - should write to log or other reaction
			return
		}

		w.Write(response)
	}
}

type JsonInterfaceHandler interface {
	Handle(w http.ResponseWriter, req *http.Request) (interface{}, error)
}

func JsonHandler(factory func(*http.Request) JsonInterfaceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")

		d, err := factory(req).Handle(w, req)
		if err != nil {
			// TODO: handle error - should write to log or other reaction
			return
		}

		response, err := json.Marshal(d)
		if err != nil {
			// TODO: handle error - should write to log or other reaction
			return
		}

		w.Write(response)
	}
}

func JsonReadUnmarshal(resp *http.Response, data interface{}) (err error) {
    contentEncoding := ""

    if encoding := resp.Header["Content-Encoding"]; encoding != nil && len(encoding) > 0 {
        contentEncoding = encoding[0]
    }

    response, err := CompressReadAll(contentEncoding, resp.Body)
    if err != nil {
        err = errors.Wrap(err, "failed to read the response")
        return
    }

    if err = json.Unmarshal(response, data); err != nil {
        err = errors.Wrap(err, "failed to unmarshal JSON")
        return
    }

    return
}
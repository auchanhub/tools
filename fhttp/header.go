package fhttp

import (
	"mime"
	"strings"
	"os"
	"net/http"
	"github.com/pkg/errors"
	"time"
)

func HeaderCharsetGet(pageUrl string, header http.Header) (charset string, err error) {
	content_type := header.Get("Content-Type")
	if content_type == "" {
		err = errors.Wrapf(os.ErrInvalid, "failed to parse charset for '", pageUrl, "', the header is not exists")
		return
	}

	mediaType, params, err := mime.ParseMediaType(header.Get("Content-Type"))
	if !strings.HasPrefix(mediaType, "text/") {
		err = errors.Wrapf(err, "failed to parse charset for '", pageUrl, "', Content-Type '", content_type, "' is wrong")
		return
	}

	var (
		ok bool
	)

	if charset, ok = params["charset"]; !ok {
		charset = "utf-8"
	}

	return
}

func HeaderDateGet(pageUrl string, header http.Header) (dateModify time.Time, err error) {
	date := header.Get("date")
	if date == "" {
		err = errors.Wrapf(os.ErrInvalid, "failed to parse Date for '", pageUrl, "', the header is not exists")
		return
	}

	dateModify, err = time.Parse(time.RFC1123, date)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse Date for '", pageUrl, "', parser error")
		return
	}

	return
}

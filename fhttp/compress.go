package fhttp

import (
	"compress/gzip"
	"net/http"
	"io"
	"strings"
	"compress/flate"
	"io/ioutil"
	"github.com/pkg/errors"
)

// The http handler is for compress response with choose type of compress by request headers
const (
	compressGzip = 0x1
	compressFlate = 0x2
	compressAll = compressGzip | compressFlate
)

type flushWriter interface {
	Write([]byte) (int, error)
	Flush() error
}

type compressWriter struct {
	mode     int
	response http.ResponseWriter
	writer   io.Writer
	flusher  flushWriter
}

func compressResponseWriter(w http.ResponseWriter, req *http.Request) (response http.ResponseWriter, flusher flushWriter) {
	acceptEncoding := req.Header.Get("Accept-Encoding")
	acceptCompresser := 0

	for i, method := range strings.Fields(acceptEncoding) {
		switch strings.TrimSuffix(method, ",") {
		case "gzip":
			acceptCompresser |= compressGzip

		case "deflate":
			acceptCompresser |= compressFlate
		}

		if acceptCompresser == compressAll || i >= 20 {
			break
		}
	}

	if acceptCompresser & compressAll != 0 {
		var (
			mode int
			writer io.Writer
			err error
		)

		switch {
		case acceptCompresser & compressFlate != 0:
			if flusher, err = flate.NewWriter(w, gzip.BestSpeed); err == nil && flusher != nil {
				w.Header().Set("Content-Encoding", "deflate")

				mode = compressGzip

				writer = flusher
			} else {
				// TODO: catch compresser error
				writer = w
			}

		case acceptCompresser & compressGzip != 0:
			w.Header().Set("Content-Encoding", "gzip")

			mode = compressFlate

			flusher = gzip.NewWriter(w)
			writer = flusher
		}

		response = &compressWriter{
			mode: mode,
			response: w,
			writer:   writer,
			flusher:  flusher,
		}
	} else {
		response = w
	}

	return
}

func (w *compressWriter) Header() http.Header {
	return w.response.Header()
}

func (w *compressWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *compressWriter) WriteHeader(header int) {
	w.response.WriteHeader(header)
}

func CompressHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		response, flusher := compressResponseWriter(w, req)

		handler(response, req)

		if flusher != nil {
			flusher.Flush()
		}
	}
}

func CompressReadAll(w http.ResponseWriter, reader io.Reader) (body []byte, err error) {
	var (
		uncompress io.Reader
	)

	switch w.Header().Get("Content-Encoding") {
	case "gzip":
		uncompress, err = gzip.NewReader(reader)

		if err != nil {
			err = errors.Wrap(err, "failed to create gzip reader for read and uncompress data")
			return
		}

	case "deflate":
		uncompress = flate.NewReader(reader)

	default:
		uncompress = reader
	}

	if body, err = ioutil.ReadAll(uncompress); err != nil && err != io.ErrUnexpectedEOF {
		err = errors.Wrap(err, "failed to read all data from the uncompress reader")
		return
	}

	if err == io.ErrUnexpectedEOF {
		err = nil
	}

	return
}

package fhttp

import (
	"testing"
	"net/http/httptest"
	"net/http"

	tools "../../"
	"../../fhttp"
	"reflect"
)

func testCompressHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello world!"))
}

type testCompressSuite struct {
	accept string
	check  func(string) bool
}

func (o *testCompressSuite) test(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/test", nil)

	if o.accept != "" {
		req.Header.Set("Accept-Encoding", o.accept)
	}

	w := httptest.NewRecorder()

	http.HandlerFunc(fhttp.CompressHandler(testCompressHandler)).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("the status of the response is", w.Code, ", but the expected status is", http.StatusOK)
		return
	}

	if encoding := w.Header().Get("Content-Encoding"); !o.check(encoding) {
		t.Error("the encoding of the response is", encoding, ", but the expected encoding is 'gzip' or 'deflate'")
		return
	}

	expectResut := "Hello world!"

	if result, err := fhttp.CompressReadAll(w, w.Body); err != nil || len(result) == 0 {
		t.Error("failed to read the response", tools.ErrorsDump(err))
	} else if !reflect.DeepEqual(string(result), expectResut) {
		t.Error("the response is", string(result), ", but the expected result should contains", expectResut)
	}
}

func TestCompressResponseWriterAll(t *testing.T) {
	testData := []testCompressSuite{
		{
			accept:"gzip, deflate, sdch",
			check: func(encoding string) bool {
				return encoding == "gzip" || encoding == "deflate"
			},
		},
		{
			accept:"gzip",
			check: func(encoding string) bool {
				return encoding == "gzip"
			},
		},
		{
			accept:"deflate",
			check: func(encoding string) bool {
				return encoding == "deflate"
			},
		},
		{
			accept:"sdch, unknown",
			check: func(encoding string) bool {
				return encoding == ""
			},
		},
		{
			accept:"",
			check: func(encoding string) bool {
				return encoding == ""
			},
		},
	}

	for _, test := range testData {
		test.test(t)
	}

}

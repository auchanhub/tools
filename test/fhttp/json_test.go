package fhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"reflect"

	"../../fhttp"
	tools "../../"
)

type testJsonData struct {
	Accept   string
	Encoding int
}

func testJsonHandler(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return &testJsonData{
		Accept: "encoding",
		Encoding: 123,
	}, nil
}

func (o *testJsonData) Handle(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	o.Accept = "encoding"
	o.Encoding = 123

	return o, nil
}

func TestJsonFuncHandle(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/test", nil)

	w := httptest.NewRecorder()

	http.HandlerFunc(fhttp.JsonFuncHandle(testJsonHandler)).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("the status of the response is", w.Code, ", but the expected status is", http.StatusOK)
		return
	}

	expectResult := `{"Accept":"encoding","Encoding":123}`

	if result, err := fhttp.CompressReadAll(w, w.Body); err != nil || len(result) == 0 {
		t.Error("failed to read the response", tools.ErrorsDump(err))
	} else if !reflect.DeepEqual(string(result), expectResult) {
		t.Error("the response is", string(result), ", but the expected result should contains", expectResult)
	}
}

func TestJsonInterfaceHandle(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/test", nil)

	w := httptest.NewRecorder()

	http.HandlerFunc(fhttp.JsonHandler(func(*http.Request) fhttp.JsonInterfaceHandler {
		return &testJsonData{}
	})).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("the status of the response is", w.Code, ", but the expected status is", http.StatusOK)
		return
	}

	expectResult := `{"Accept":"encoding","Encoding":123}`

	if result, err := fhttp.CompressReadAll(w, w.Body); err != nil || len(result) == 0 {
		t.Error("failed to read the response", tools.ErrorsDump(err))
	} else if !reflect.DeepEqual(string(result), expectResult) {
		t.Error("the response is", string(result), ", but the expected result should contains", expectResult)
	}
}
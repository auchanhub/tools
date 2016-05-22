package fhttp

import (
	"testing"
	"reflect"
	"strconv"

	"net/http"
	"net/http/httptest"

	"../../fhttp"
	tools "../../"
)

type testChainData struct {
	Accept   string
	Encoding int
}

func (o *testChainData) Handle(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return o, nil
}

func TestChainHandler(t *testing.T) {
	handler := fhttp.CompressHandler(
		fhttp.JsonInterfaceHandle(
			func(req *http.Request) fhttp.JsonInterfaceHandler {
				return fhttp.BindHandlerParams(req, &testChainData{}).(fhttp.JsonInterfaceHandler)
			}))

	accept := "encoding"
	encoding := 123

	req, _ := http.NewRequest("GET", "/api/test", nil)

	values := req.URL.Query()
	values.Add("Accept", accept)
	values.Add("Encoding", strconv.Itoa(encoding))

	req.URL.RawQuery = values.Encode()

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

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
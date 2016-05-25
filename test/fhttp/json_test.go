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

    if result, err := fhttp.CompressReadAll(w.Header().Get("Content-Encoding"), w.Body); err != nil || len(result) == 0 {
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

    if result, err := fhttp.CompressReadAll(w.Header().Get("Content-Encoding"), w.Body); err != nil || len(result) == 0 {
        t.Error("failed to read the response", tools.ErrorsDump(err))
    } else if !reflect.DeepEqual(string(result), expectResult) {
        t.Error("the response is", string(result), ", but the expected result should contains", expectResult)
    }
}

func TestJsonReadUnmarshal(t *testing.T) {
    port, test_server, err := fhttp.TestServer(http.HandlerFunc(fhttp.JsonHandler(func(*http.Request) fhttp.JsonInterfaceHandler {
        return &testJsonData{}
    })))
    if err != nil {
        t.Error("failed to create test server", tools.ErrorsDump(err))
        return
    }
    test_server.Start()
    defer test_server.Close()

    req, _ := http.NewRequest("GET", "/api/test", nil)

    resp, err := fhttp.NewPool("localhost", port).Do(req)
    if err != nil {
        t.Error("failed to execute test request", tools.ErrorsDump(err))
        return
    }

    expectResult := &testJsonData{
        Accept:   "encoding",
        Encoding: 123,
    }

    existResult := &testJsonData{}

    err = fhttp.JsonReadUnmarshal(resp, existResult)
    if err != nil {
        t.Error("failed to unmarshal a response for the test request", tools.ErrorsDump(err))
        return
    }

    if !reflect.DeepEqual(existResult, expectResult) {
        t.Error("the response is", existResult, ", but the expected result should contains", expectResult)
    }
}
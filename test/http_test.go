package test

import (
	"testing"
	"net/http"

	tools "../"
	"reflect"
	"net/url"
)

func TestRequestParseParamsGET(t *testing.T) {
	type testData struct {
		Help string `schema:"help"`
		Params int `schema:"params"`
	}

	expectData := &testData{
		Help: "me",
		Params: 123,
	}

	req, _ := http.NewRequest("GET", "/api/test", nil)

	values := req.URL.Query()
	values.Add("help", "me")
	values.Add("params", "123")

	req.URL.RawQuery = values.Encode()

	existData := &testData{}

	if err:=tools.RequestParseParams(req, existData); err!=nil {
		t.Error("failed to parse request params", tools.ErrorsDump(err))
	}

	if !reflect.DeepEqual(expectData, existData) {
		t.Error("the result of request parse is\n", existData, "\n, but the expected result is\n", expectData)
	}
}

func TestRequestParseParamsPOST(t *testing.T) {
	type testData struct {
		Help string `schema:"help"`
		Params int `schema:"params"`
	}

	expectData := &testData{
		Help: "me",
		Params: 123,
	}

	req, _ := http.NewRequest("POST", "/api/test", nil)
	req.PostForm = make(url.Values)
	req.PostForm.Set("help", "me")
	req.PostForm.Set("params", "123")

	existData := &testData{}

	if err:=tools.RequestParseParams(req, existData); err!=nil {
		t.Error("failed to parse request params", tools.ErrorsDump(err))
	}

	if !reflect.DeepEqual(expectData, existData) {
		t.Error("the result of request parse is\n", existData, "\n, but the expected result is\n", expectData)
	}
}
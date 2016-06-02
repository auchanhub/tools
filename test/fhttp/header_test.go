package fhttp

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"../../fhttp"
	"time"
	"reflect"
)

func TestHeaderCharsetGetParams(t *testing.T) {
	if charset, err := fhttp.HeaderCharsetGet("test", nil); charset != "" || err == nil {
		t.Error("failed to check the parmeters. The error is nil, but should be set")
		return
	}

	pageUrl := "/api/test"

	req, _ := http.NewRequest("GET", pageUrl, nil)

	w := httptest.NewRecorder()

	http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("X-Content-Type", "test")
		w.Write([]byte("Hello world!"))
	}).ServeHTTP(w, req)

	if charset, err := fhttp.HeaderCharsetGet(pageUrl, w.Header()); charset != "utf-8" || err != nil {
		t.Error("failed to check the parmeters. The error is nil, but should be set", charset, err)
		return
	}
}

func TestHeaderCharsetGet(t *testing.T) {
	pageUrl := "/api/test"

	req, _ := http.NewRequest("GET", pageUrl, nil)

	w := httptest.NewRecorder()

	http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "text/json; charset=windows-1251")
		w.Write([]byte("Hello world!"))
	}).ServeHTTP(w, req)

	if charset, err := fhttp.HeaderCharsetGet(pageUrl, w.Header()); charset != "windows-1251" || err != nil {
		t.Error("failed to parse the charset. It is ", charset, ", but it should be windows-1251", err)
		return
	}
}

func TestHeaderDateGet(t *testing.T) {
	pageUrl := "/api/test"

	req, _ := http.NewRequest("GET", pageUrl, nil)

	now := time.Now()

	w := httptest.NewRecorder()

	http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Date", now.Format(time.RFC1123))
		w.Write([]byte("Hello world!"))
	}).ServeHTTP(w, req)

	if date, err := fhttp.HeaderDateGet(pageUrl, w.Header()); !reflect.DeepEqual(date.Format(time.RFC1123), now.Format(time.RFC1123)) || err != nil {
		t.Error("failed to parse the date. It is ", date.Format(time.RFC1123), ", but it should be", now.Format(time.RFC1123), err)
		return
	}
}

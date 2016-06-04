package fhttp

import (
	"testing"
	"net/http"
	"fmt"
	"reflect"

	"../../fhttp"
	tools "../../"
)

func TestNewPool(t *testing.T) {
	pool := fhttp.NewPool("localhost", 9000)
	client := pool.Get()
	client2 := pool.Get()

	pool.Put(client)

	client = pool.Get()

	if client == nil || client2 == nil {
		t.Error("failed to create client from pool")
	}
}

func TestPoolRequest(t *testing.T) {
	testHost := "www.osp.ru"

	test_server, err := fhttp.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, ", r.Host)
	})
	if err != nil {
		t.Error("failed to create a test http server and jstart listen a port", test_server.Port, tools.ErrorsDump(err))
		return
	}

	test_server.Start()
	defer test_server.Close()

	req, err := http.NewRequest("GET", fmt.Sprint("https://", testHost, "/api/test"), nil)
	if err != nil {
		t.Error("failed to create the new request", err)
		return
	}

	pool := fhttp.NewPool("localhost", test_server.PortTLS).SetSkipCertVerify(true)

	resp, err := pool.Do(req)
	if err != nil {
		t.Error("failed to execute the new request", err)
		return
	}

	expectResult := "Hello, " + testHost

	if result, err := fhttp.CompressReadAll(resp.Header.Get("Content-Encoding"), resp.Body); err != nil || len(result) == 0 {
		t.Error("failed to read the response", tools.ErrorsDump(err))
	} else if !reflect.DeepEqual(string(result), expectResult) {
		t.Error("the response is", string(result), ", but the expected result should contains", expectResult)
	}
}

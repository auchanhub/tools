package fhttp

import (
    "fmt"
    "net"
    "net/http"
    "math/rand"
    "net/http/httptest"

    "github.com/pkg/errors"
)

func TestServer(handler http.HandlerFunc) (int, *httptest.Server, error) {
    var (
        listener net.Listener
        err error
        address string
        port int
    )

    // try three times to listen a port
    for err_count := 0; err_count < 3; err_count++ {
        port = 10000 + int(rand.Int31n(50000))
        address = fmt.Sprintf("localhost:%d", port)

        listener, err = net.Listen("tcp", address)
        if err == nil {
            break
        }
    }

    if err != nil {
        return 0, nil, errors.Wrapf(err, "failed to start listen a port", address, err)
    }

    return port, &httptest.Server{
        Listener: listener,
        Config:   &http.Server{
            Handler: handler,
        },
    }, nil
}

package fhttp

import (
    "net/http"
    "sync"
    "net"
    "time"
    "github.com/pkg/errors"
    "os"
    "fmt"
)

type Pool struct {
    sync.Pool

    host    string
    address string
}

func NewPool(address string, port int) (pool *Pool) {
    pool = &Pool{
        host: address,
        address: fmt.Sprintf("%s:%d", address, port),
    }

    pool.New = func() interface{} {
        return &http.Client{
            Timeout: 5 * time.Second,

            Transport: &http.Transport{
                Dial: (&net.Dialer{
                    Timeout: 5 * time.Second,
                }).Dial,

                TLSHandshakeTimeout: 5 * time.Second,
            },
        }
    }

    return
}

func (p *Pool) Get() *http.Client {
    return p.Pool.Get().(*http.Client)
}

func (p *Pool) Put(client *http.Client) {
    p.Pool.Put(client)
}

func (p *Pool) Do(req *http.Request) (*http.Response, error) {
    client := p.Get()
    defer p.Put(client)

    if client == nil {
        return nil, errors.Wrapf(os.ErrNotExist, "failed to get a connection to '%v'", p.address)
    }

    req.Host = p.address
    req.URL.Host = p.address
    req.URL.Scheme = "http"

    return client.Do(req)
}

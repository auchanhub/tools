package fhttp

import (
	"net/http"
	"sync"
	"net"
	"time"
	"github.com/pkg/errors"
	"os"
	"fmt"
	"crypto/tls"
)

type RequestExecutor interface {
	Do(req *http.Request) (*http.Response, error)
}

type Executor struct {
	SkipCertVerify bool
}

func (o *Executor) Do(req *http.Request) (*http.Response, error) {
	skipCertVerify := true
	if o != nil {
		skipCertVerify = o.SkipCertVerify
	}

	return NewClient(skipCertVerify).Do(req)
}

func NewClient(skipCertVerify bool) *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,

		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 7 * time.Second,
			}).Dial,

			TLSHandshakeTimeout: 5 * time.Second,

			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipCertVerify},
		},
	}
}


type Pool struct {
	sync.Pool

	host           string
	address        string
	skipCertVerify bool
}

func NewPool(address string, port int) (pool *Pool) {
	pool = &Pool{
		host: address,
		address: fmt.Sprintf("%s:%d", address, port),
	}

	pool.Pool.New = func() interface{} {
		return NewClient(pool.skipCertVerify)
	}

	return
}

func (o *Pool) SetSkipCertVerify(skip bool) *Pool {
	o.skipCertVerify = true

	return o
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

	req.URL.Host = p.address

	return client.Do(req)
}

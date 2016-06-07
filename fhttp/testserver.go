package fhttp

import (
	"fmt"
	"net"
	"net/http"
	"math/rand"
	"net/http/httptest"

	"github.com/pkg/errors"
	"crypto/tls"
)

// From go source of the test server
// LocalhostCert is a PEM-encoded TLS cert with SAN IPs
// "127.0.0.1" and "[::1]", expiring at Jan 29 16:00:00 2084 GMT.
// generated from src/crypto/tls:
// go run generate_cert.go  --rsa-bits 1024 --host 127.0.0.1,::1,example.com --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
var localhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIICEzCCAXygAwIBAgIQMIMChMLGrR+QvmQvpwAU6zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9SjY1bIw4
iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZBl2+XsDul
rKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQABo2gwZjAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAuBgNVHREEJzAlggtleGFtcGxlLmNvbYcEfwAAAYcQAAAAAAAAAAAAAAAA
AAAAATANBgkqhkiG9w0BAQsFAAOBgQCEcetwO59EWk7WiJsG4x8SY+UIAA+flUI9
tyC4lNhbcF2Idq9greZwbYCqTTTr2XiRNSMLCOjKyI7ukPoPjo16ocHj+P3vZGfs
h1fIw3cSS2OolhloGw/XM6RWPWtPAlGykKLciQrBru5NAPvCMsb/I1DAceTiotQM
fblo6RBxUQ==
-----END CERTIFICATE-----`)

// LocalhostKey is the private key for localhostCert.
var localhostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9
SjY1bIw4iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZB
l2+XsDulrKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQAB
AoGAGRzwwir7XvBOAy5tM/uV6e+Zf6anZzus1s1Y1ClbjbE6HXbnWWF/wbZGOpet
3Zm4vD6MXc7jpTLryzTQIvVdfQbRc6+MUVeLKwZatTXtdZrhu+Jk7hx0nTPy8Jcb
uJqFk541aEw+mMogY/xEcfbWd6IOkp+4xqjlFLBEDytgbIECQQDvH/E6nk+hgN4H
qzzVtxxr397vWrjrIgPbJpQvBsafG7b0dA4AFjwVbFLmQcj2PprIMmPcQrooz8vp
jy4SHEg1AkEA/v13/5M47K9vCxmb8QeD/asydfsgS5TeuNi8DoUBEmiSJwma7FXY
fFUtxuvL7XvjwjN5B30pNEbc6Iuyt7y4MQJBAIt21su4b3sjXNueLKH85Q+phy2U
fQtuUE9txblTu14q3N7gHRZB4ZMhFYyDy8CKrN2cPg/Fvyt0Xlp/DoCzjA0CQQDU
y2ptGsuSmgUtWj3NM9xuwYPm+Z/F84K6+ARYiZ6PYj013sovGKUFfYAqVXVlxtIX
qyUBnu3X9ps8ZfjLZO7BAkEAlT4R5Yl6cGhaJQYZHOde3JEMhNRcVFMO8dJDaFeo
f9Oeos0UUothgiDktdQHxdNEwLjQf7lJJBzV+5OtwswCWA==
-----END RSA PRIVATE KEY-----`)

type TestServer struct {
	Port      int
	PortTLS   int

	Pool      *Pool
	PoolTLS   *Pool

	server    *httptest.Server
	serverTLS *httptest.Server
}

func NewTestServer(handler http.HandlerFunc) (test_server *TestServer, err error) {
	cert, _ := tls.X509KeyPair(localhostCert, localhostKey)

	test_server = &TestServer{
		server: &httptest.Server{
			Config:   &http.Server{
				Handler: handler,
			},
		},
		serverTLS: &httptest.Server{
			TLS: &tls.Config{
				NextProtos: []string{"http/1.1"},
				Certificates: []tls.Certificate{cert},
			},
			Config:   &http.Server{
				Handler: handler,
			},
		},
	}

	test_server.Port, test_server.server.Listener, err = test_server.genListener()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to init listen a port", test_server.Port, err)
	}
	test_server.Pool = NewPool("localhost", test_server.Port).SetSkipCertVerify(true)

	test_server.PortTLS, test_server.serverTLS.Listener, err = test_server.genListener()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to init TLS listen a port", test_server.PortTLS, err)
	}
	test_server.PoolTLS = NewPool("localhost", test_server.PortTLS).SetSkipCertVerify(true)

	return
}

func (o *TestServer) genListener() (port int, listener net.Listener, err error) {
	var (
		address string
	)
	// try three times to listen a port
	for err_count := 0; err_count < 10; err_count++ {
		port = 10000 + int(rand.Int31n(50000))
		address = fmt.Sprintf("localhost:%d", port)

		listener, err = net.Listen("tcp", address)
		if err == nil {
			break
		}
	}

	return
}

func (o *TestServer) Start() {
	o.server.Start()
	o.serverTLS.StartTLS()
}

func (o *TestServer) Do(req *http.Request) (*http.Response, error) {
	switch req.URL.Scheme {
	case "https":
		return o.PoolTLS.Do(req)
	default:
		return o.Pool.Do(req)
	}
}

func (o *TestServer) Close() {
	o.server.Close()
	o.serverTLS.Close()
}

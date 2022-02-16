package fetcher

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDefaultClient(t *testing.T) {
	c := &Client{}
	d, err := c.CreateDoer()
	if err != nil {
		t.Fatal(err)
	}
	client := d.(*http.Client)
	if client == nil {
		t.Fatal()
	}
	if client.Timeout != DefaultTimeout {
		t.Fatal(client)
	}
	transport := client.Transport.(*http.Transport)
	if transport == nil {
		t.Fatal()
	}
	if transport.MaxIdleConns != DefaultMaxIdleConns ||
		transport.TLSHandshakeTimeout != DefaultTLSHandshakeTimeout ||
		transport.IdleConnTimeout != DefaultIdleConnTimeout {
		t.Fatal(transport)
	}
}

func TestClientConfig(t *testing.T) {
	c := &Client{
		TimeoutInSecond:             15,
		TLSHandshakeTimeoutInSecond: 16,
		IdleConnTimeoutInSecond:     17,
		MaxIdleConns:                18,
	}
	d, err := c.CreateDoer()
	if err != nil {
		t.Fatal(err)
	}
	client := d.(*http.Client)
	if client == nil {
		t.Fatal()
	}
	if client.Timeout != 15*time.Second {
		t.Fatal(client)
	}
	transport := client.Transport.(*http.Transport)
	if transport == nil {
		t.Fatal()
	}
	if transport.MaxIdleConns != 18 ||
		transport.TLSHandshakeTimeout != 16*time.Second ||
		transport.IdleConnTimeout != 17*time.Second {
		t.Fatal(transport)
	}

}

func TestProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("proxied"))
		if err != nil {
			panic(err)
		}
	}))
	defer proxy.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("notproxied"))
		if err != nil {
			panic(err)
		}
	}))
	defer server.Close()

	client := Client{
		Proxy: proxy.URL,
	}
	err := client.SelfCheck()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if string(bs) != "proxied" {
		t.Error(string(bs))
	}
}

func TestNoProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("proxied"))
		if err != nil {
			panic(err)
		}
	}))
	defer proxy.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("notproxied"))
		if err != nil {
			panic(err)
		}
	}))
	defer server.Close()

	clients := Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := clients.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if string(bs) != "notproxied" {
		t.Error(string(bs))
	}

}

func TestGzip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		writer := gzip.NewWriter(buf)
		writer.Write([]byte("testtest"))
		writer.Close()
		w.Header().Add("Content-Encoding", "gzip")
		w.Write(buf.Bytes())
	}))
	defer server.Close()

	client := Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if string(bs) != "testtest" {
		t.Error(string(bs))
	}

}

func TestNilClient(t *testing.T) {
	s := newEchoServer()
	defer s.Close()
	p := MustPreset(&ServerInfo{
		URL: s.URL,
	})
	resp, err := p.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestClientExec(t *testing.T) {
	f := New()
	doer := &Client{}
	err := doer.Exec(f)
	if err != nil {
		panic(err)
	}
	if f.Doer != doer {
		t.Fatal(f)
	}
}

func TestClientClone(t *testing.T) {
	c := &Client{}
	cloned := c.Clone()
	c.Proxy = "1234"
	if cloned.Proxy == c.Proxy {
		t.Fatal(cloned)
	}
}

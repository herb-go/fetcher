package fetcher

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

//DefaultDoer default doer used to do http request if not setted.
var DefaultDoer = func() Doer {
	return http.DefaultClient
}

//DefaultTimeout default client timeout
var DefaultTimeout = 120 * time.Second

//DefaultMaxIdleConns default client max idle conns
var DefaultMaxIdleConns = 20

//DefaultIdleConnTimeout default client idle conn timeout.
var DefaultIdleConnTimeout = 120 * time.Second

//DefaultTLSHandshakeTimeout default client tls handshake time out.
var DefaultTLSHandshakeTimeout = 30 * time.Second

//Client http client config struct
//Value should not changed after first "DO" call.
type Client struct {
	//TimeoutInSecond timeout in secound
	//Default value is 120.
	TimeoutInSecond int64
	//MaxIdleConns max idel conns.
	//Default value is 20
	MaxIdleConns int
	//IdleConnTimeoutInSecond idel conn timeout in second.
	//Default value is 120
	IdleConnTimeoutInSecond int64
	//TLSHandshakeTimeoutInSecond tls handshake timeout in secound.
	//default value is 30.
	TLSHandshakeTimeoutInSecond int64
	//Proxy proxy url.
	//If set to empty string,clients will not use proxy.
	//Default value is empty string.
	Proxy  string
	locker sync.Mutex
	doer   Doer
}

//CreateDoer create doer.
//Return doer createrd and any error if raised.
func (c *Client) CreateDoer() (Doer, error) {
	var timeout time.Duration
	if c.TimeoutInSecond > 0 {
		timeout = time.Duration(c.TimeoutInSecond) * time.Second
	} else {
		timeout = DefaultTimeout
	}
	transport := c.getTransport()
	transport.Proxy = URLToProxy(c.Proxy)
	client := http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	return &client, nil
}

func (c *Client) getDoer() (Doer, error) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.doer != nil {
		return c.doer, nil
	}
	d, err := c.CreateDoer()
	if err != nil {
		return nil, err
	}
	return d, nil
}

//Do create request ,client with given f.etcher and commands and fetch
//Return http response and any error if raised.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	d, err := c.getDoer()
	if err != nil {
		return nil, err
	}
	return d.Do(req)
}

//SelfCheck self check.
//Return any error if raised.
func (c *Client) SelfCheck() error {
	_, err := c.getDoer()
	return err
}

//URLToProxy Convert url to fixed proxy.
//Return proxy created and any error if raised.
func URLToProxy(index string) func(*http.Request) (*url.URL, error) {
	if index == "" {
		return http.ProxyFromEnvironment
	}
	url, err := url.Parse(index)
	if err != nil {
		return http.ProxyFromEnvironment
	}
	return http.ProxyURL(url)
}

func (c *Client) getTransport() *http.Transport {
	var maxIdleCoons = c.MaxIdleConns
	if maxIdleCoons == 0 {
		maxIdleCoons = DefaultMaxIdleConns
	}
	var idleConnTimeout time.Duration
	if c.IdleConnTimeoutInSecond <= 0 {
		idleConnTimeout = DefaultIdleConnTimeout
	} else {
		idleConnTimeout = time.Duration(c.IdleConnTimeoutInSecond) * time.Second
	}
	var tlsHandshakeTimeout time.Duration
	if c.TLSHandshakeTimeoutInSecond <= 0 {
		tlsHandshakeTimeout = DefaultTLSHandshakeTimeout
	} else {
		tlsHandshakeTimeout = time.Duration(c.TLSHandshakeTimeoutInSecond) * time.Second
	}

	return &http.Transport{
		Proxy:               URLToProxy(c.Proxy),
		MaxIdleConns:        maxIdleCoons,
		IdleConnTimeout:     idleConnTimeout,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}
}

//DoerFactory doer factory
type DoerFactory interface {
	//CreateDoer create doer.
	//Return doer createrd and any error if raised.
	CreateDoer() (*Doer, error)
}

// Doer doer interface
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

//Do create request ,client with given f.etcher and commands and fetch
//Return http response and any error if raised.
func Do(f *Fetcher, b ...Command) (*Response, error) {

	err := Exec(f, b...)
	if err != nil {
		return nil, err
	}
	resp, err := f.Fetch()
	if err != nil {
		return nil, err
	}
	return ConvertResponse(resp), nil
}

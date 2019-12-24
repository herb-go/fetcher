package fetcher

import (
	"io"
	"net/http"
	"net/url"
)

//Fetcher http request fetcher struct.
//New fetcher should be created when new http request buildding.
//You should not edit Fetcher value directly,use Command and Preset instead.
type Fetcher struct {
	//URL http url used to create http request
	URL *url.URL
	//Header http header used to create http request
	Header http.Header
	//Method http method used to create http request
	Method string
	//Body request body
	Body io.Reader
	//Builders request builder which should called in order after http request created.
	Builders []func(*http.Request) error
	//Doer http client by which will do request
	Doer Doer
}

//AppendBuilder append request builders to fetcher.
//Fetcher builders will be cloned.
func (f *Fetcher) AppendBuilder(b ...func(*http.Request) error) {
	f.Builders = append(CloneRequestBuilders(f.Builders), b...)
}

//Raw create raw http request,doer and any error if raised.
func (f *Fetcher) Raw() (*http.Request, Doer, error) {
	url := f.URL.String()
	req, err := http.NewRequest(f.Method, url, f.Body)
	if err != nil {
		return nil, nil, err
	}
	MergeHeader(req.Header, f.Header)
	if f.Doer == nil {
		return req, DefaultDoer(), nil
	}
	return req, f.Doer, nil
}

//Fetch create http requuest and fetch.
//Return http response and any error if raised.
func (f *Fetcher) Fetch() (*http.Response, error) {
	req, doer, err := f.Raw()
	if err != nil {
		return nil, err
	}
	return doer.Do(req)
}

//New create new fetcher
func New() *Fetcher {
	return &Fetcher{
		Builders: []func(*http.Request) error{},
	}
}

//Fetch create new fetcher ,exec commands and fetch response.
//Return http response and any error if raised.
func Fetch(cmds ...Command) (*Response, error) {
	f := New()
	return Do(f, cmds...)
}

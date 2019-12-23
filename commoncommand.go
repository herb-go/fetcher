package fetcher

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

//URL create new command which modify fetcher url to given url
func URL(u *url.URL) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.URL = u
		return nil
	})
}

//Method command which modify fetcher method to given method
type Method string

//Exec exec command to modify fetcher.
//Return any error if raised.
func (m Method) Exec(f *Fetcher) error {
	f.Method = string(m)
	return nil
}

//PathPrefix command which modify fetcher url with given path prefix
type PathPrefix string

//Exec exec command to modify fetcher.
//Return any error if raised.
func (p PathPrefix) Exec(f *Fetcher) error {
	f.URL.Path = f.URL.Path + string(p)
	return nil
}

//Replace command which modify fetcher path by given placeholder and value.
func Replace(placeholder string, value string) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.URL.Path = strings.NewReplacer(placeholder, value).Replace(f.URL.Path)
		return nil
	})
}

var (
	//Post http POST method Command.
	Post = Method("POST")
	//Get http GET method Command.
	Get = Method("GET")
	//Put http PUT method Command.
	Put = Method("PUT")
	//Delete http DELETE method Command.
	Delete = Method("DELETE")
)

//Body command which modify fetcher body to given reader.
func Body(body io.Reader) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.Body = body
		return nil
	})
}

//Header command which merge fetcher header by given reader.
func Header(h http.Header) Command {
	return CommandFunc(func(f *Fetcher) error {
		MergeHeader(f.Header, h)
		return nil
	})
}

//SetDoer command which modify fetcher doer to given doer.
func SetDoer(d Doer) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.Doer = d
		return nil
	})
}

//RequestBuilderProvider request builder provider interface.
type RequestBuilderProvider interface {
	BuildRequest(*http.Request) error
}

//RequestBuilder command which append given request builder to fetcher.
func RequestBuilder(p RequestBuilderProvider) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.AppendBuilder(p.BuildRequest)
		return nil
	})
}

//HeaderBuilderProvier header builde provider
type HeaderBuilderProvier interface {
	BuildHeader(http.Header) error
}

//HeaderBuilder command which modify fetcher header by given header builder provider.
func HeaderBuilder(p HeaderBuilderProvier) Command {
	return CommandFunc(func(f *Fetcher) error {
		return p.BuildHeader(f.Header)
	})
}

//MethodBuilderProvider method builder provider
type MethodBuilderProvider interface {
	RequestMethod() (string, error)
}

//MethodBuilder command which modify fetcher method by given method builder provider.
func MethodBuilder(p MethodBuilderProvider) Command {
	return CommandFunc(func(f *Fetcher) error {
		m, err := p.RequestMethod()
		if err != nil {
			return err
		}
		f.Method = m
		return nil
	})
}

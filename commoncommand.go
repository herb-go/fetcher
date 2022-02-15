package fetcher

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

//ParsedURL create new command which modify fetcher url to given prased url
func ParsedURL(u *url.URL) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.URL = u
		return nil
	})
}

//URL create new command which modify fetcher url to given url
func URL(u string) Command {
	return CommandFunc(func(f *Fetcher) error {
		furl, err := url.Parse(u)
		if err != nil {
			return err
		}
		f.URL = furl
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

//PathPrefix command which modify fetcher url with given path prefix
type PathPrefix string

//Exec exec command to modify fetcher.
//Return any error if raised.
func (p PathPrefix) Exec(f *Fetcher) error {
	f.URL.Path = string(p) + f.URL.Path
	return nil
}

//PathSuffix command which modify fetcher url with given path suffix
type PathSuffix string

//Exec exec command to modify fetcher.
//Return any error if raised.
func (p PathSuffix) Exec(f *Fetcher) error {
	f.URL.Path = f.URL.Path + string(p)
	return nil
}

//Host command which modify fetcher url with given host
type Host string

//Exec exec command to modify fetcher.
//Return any error if raised.
func (h Host) Exec(f *Fetcher) error {
	f.URL.Host = string(h)
	return nil
}

//Replace command which modify fetcher path by given placeholder and value.
func Replace(placeholder string, value string) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.URL.Path = strings.NewReplacer(placeholder, value).Replace(f.URL.Path)
		return nil
	})
}

//Body command which modify fetcher body to given reader.
func Body(body io.Reader) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.Body = body
		return nil
	})
}

//JSONBody command which modify fetcher body to given value as json.
//Fetcher body will set to nil if v is nil.
func JSONBody(v interface{}) Command {
	return CommandFunc(func(f *Fetcher) error {
		if v == nil {
			f.Body = nil
			return nil
		}
		bs, err := json.Marshal(v)
		if err != nil {
			return err
		}
		f.Body = bytes.NewBuffer(bs)
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

//SetHeader command which set fetcher header by given key and value.
func SetHeader(key string, value string) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.Header.Set(key, value)
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

//SetQuery command which modify fetcher to set given query.
func SetQuery(name string, value string) Command {
	return CommandFunc(func(f *Fetcher) error {
		q := f.URL.Query()
		q.Set(name, value)
		f.URL.RawQuery = q.Encode()
		return nil
	})
}

//Params command which modify fetcher to set given params.
func Params(params url.Values) Command {
	return CommandFunc(func(f *Fetcher) error {
		q := f.URL.Query()
		for key := range params {
			q.Set(key, params.Get(key))
		}
		f.URL.RawQuery = q.Encode()
		return nil
	})
}

//BasicAuth command which modify fetcher  to set given basic auth info.
func BasicAuth(username string, password string) Command {
	return CommandFunc(func(f *Fetcher) error {
		f.AppendBuilder(func(r *http.Request) error {
			r.SetBasicAuth(username, password)
			return nil
		})
		return nil
	})
}

//RequestBuilderFunc request builder func type
type RequestBuilderFunc func(*http.Request) error

//Exec exec command to modify fetcher.
//Return any error if raised.
func (b RequestBuilderFunc) Exec(f *Fetcher) error {
	f.AppendBuilder(b)
	return nil

}

//RequestBuilderProvider request builder provider interface.
type RequestBuilderProvider interface {
	BuildRequest(*http.Request) error
}

//RequestBuilder command which append given request builder to fetcher.
func RequestBuilder(p RequestBuilderProvider) Command {
	return RequestBuilderFunc(p.BuildRequest)
}

//HeaderBuilderFunc header builde func
type HeaderBuilderFunc func(http.Header) error

//Exec exec command to modify fetcher.
//Return any error if raised.
func (b HeaderBuilderFunc) Exec(f *Fetcher) error {
	return b(f.Header)
}

//HeaderBuilderProvier header builde provider
type HeaderBuilderProvier interface {
	BuildHeader(http.Header) error
}

//HeaderBuilder command which modify fetcher header by given header builder provider.
func HeaderBuilder(p HeaderBuilderProvier) Command {
	return HeaderBuilderFunc(p.BuildHeader)
}

//MethodBuilderFunc method builder func
type MethodBuilderFunc func() (string, error)

//Exec exec command to modify fetcher.
//Return any error if raised.
func (b MethodBuilderFunc) Exec(f *Fetcher) error {
	m, err := b()
	if err != nil {
		return err
	}
	f.Method = m
	return nil

}

//MethodBuilderProvider method builder provider
type MethodBuilderProvider interface {
	RequestMethod() (string, error)
}

//MethodBuilder command which modify fetcher method by given method builder provider.
func MethodBuilder(p MethodBuilderProvider) Command {
	return MethodBuilderFunc(p.RequestMethod)
}

//ParamsBuilderFunc param builde func type
type ParamsBuilderFunc func(url.Values) error

//Exec exec command to modify fetcher.
//Return any error if raised.
func (b ParamsBuilderFunc) Exec(f *Fetcher) error {
	params := f.URL.Query()
	err := b(params)
	if err != nil {
		return err
	}
	f.URL.RawQuery = params.Encode()
	return nil

}

//ParamsBuilderProvier params builde provider
type ParamsBuilderProvier interface {
	BuildParams(url.Values) error
}

//ParamsBuilder command which modify fetcher header by given params builder provider.
func ParamsBuilder(p ParamsBuilderProvier) Command {
	return ParamsBuilderFunc(p.BuildParams)
}

//MultiPartWriter multipart writer command
type MultiPartWriter struct {
	body *bytes.Buffer
	*multipart.Writer
}

//Exec exec command to modify fetcher.
//Return any error if raised.
func (w *MultiPartWriter) Exec(f *Fetcher) error {
	f.Body = w.body
	f.Header.Set("Content-Type", w.FormDataContentType())
	return nil
}

//WriteFile write file with given fieldname,filename and data
func (w *MultiPartWriter) WriteFile(fieldname, filename string, src io.Reader) error {
	writer, err := w.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, src)
	return err
}

//NewMultiPartWriter create new MultiPartWriter command
func NewMultiPartWriter() *MultiPartWriter {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	return &MultiPartWriter{
		body:   buf,
		Writer: writer,
	}
}

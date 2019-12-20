package fetcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
	bytes *[]byte
}

func NewResponse(r *http.Response) *Response {
	return &Response{
		Response: r,
	}
}

type Parser interface {
	//Parse parse response data.
	//Return Respnse and any error if raised.
	Parse(*Response, error) (*Response, error)
}

type ParserFunc func(resp *Response, err error) (*Response, error)

func (p ParserFunc) Parse(resp *Response, err error) (*Response, error) {
	return p(resp, err)
}

var BytesParser = ParserFunc(func(resp *Response, err error) (*Response, error) {
	if err != nil {
		return resp, err
	}
	if resp.bytes != nil {
		return resp, nil
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.bytes = &bs
	return resp, nil

})

func StringParser(str *string) Parser {
	return ParserFunc(func(resp *Response, err error) (*Response, error) {
		resp, err = BytesParser.Parse(resp, err)
		if err != nil {
			return resp, err
		}
		bs := *resp.bytes
		s := string(bs)
		*str = s
		return resp, nil
	})
}

func JSONParser(v interface{}) Parser {
	return ParserFunc(func(resp *Response, err error) (*Response, error) {
		resp, err = BytesParser.Parse(resp, err)
		if err != nil {
			return resp, err
		}
		bs := *resp.bytes
		err = json.Unmarshal(bs, v)
		if err != nil {
			return nil, err
		}
		return resp, err
	})
}

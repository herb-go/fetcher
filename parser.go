package fetcher

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

//Response fetch response struct
type Response struct {
	*http.Response
	bytes *[]byte
}

//BodyContent read and return body content from response.
//Response body will be closed after first read.
func (r *Response) BodyContent() ([]byte, error) {
	if r.bytes != nil {
		return *r.bytes, nil
	}
	defer r.Response.Body.Close()
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.bytes = &bs
	return *r.bytes, nil
}

//Error return response body content as error.
func (r *Response) Error() string {
	bs, err := r.BodyContent()
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

//NewResponse create new response
func NewResponse() *Response {
	return &Response{}
}

//ConvertResponse convert http response to fetch response
func ConvertResponse(resp *http.Response) *Response {
	r := NewResponse()
	r.Response = resp
	return r
}

//Parser response parser interface
type Parser interface {
	//Parse parse response data.
	//Return Response and any error if raised.
	Parse(*Response) error
}

//ParserFunc parser func type
type ParserFunc func(resp *Response) error

//Parse parse response data.
//Return Response and any error if raised.
func (p ParserFunc) Parse(resp *Response) error {
	return p(resp)
}

//ShouldSuccess parser that check if status code < 300.
//Give parser will parse responese if success or respnse will be returned as error.
func ShouldSuccess(p Parser) Parser {
	return ParserFunc(func(resp *Response) error {
		if resp.StatusCode >= 300 {
			return resp
		}
		return p.Parse(resp)
	})
}

//ShouldNoError parser that check if status code < 500.
//Give parser will parse responese if success or respnse will be returned as error.
func ShouldNoError(p Parser) Parser {
	return ParserFunc(func(resp *Response) error {
		if resp.StatusCode >= 500 {
			return resp
		}
		return p.Parse(resp)
	})
}

//AsBytes create parser which parse givn byte slice from response.
func AsBytes(bytes []byte) Parser {
	return ParserFunc(func(resp *Response) error {
		bs, err := resp.BodyContent()
		if err != nil {
			return err
		}
		bytes = make([]byte, len(bs))
		copy(bytes, bs)
		return nil
	})
}

//AsString create parser which parse givn string from response.
func AsString(str *string) Parser {
	return ParserFunc(func(resp *Response) error {
		bs, err := resp.BodyContent()
		if err != nil {
			return err
		}
		s := string(bs)
		*str = s
		return nil
	})
}

//AsJSON create parser which parse givn value from response a JSON format.
func AsJSON(v interface{}) Parser {
	return ParserFunc(func(resp *Response) error {
		bs, err := resp.BodyContent()
		if err != nil {
			return err
		}
		err = json.Unmarshal(bs, v)
		if err != nil {
			return err
		}
		return err
	})
}

//FetchAndParse fetch response and prase response with given commands and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func FetchAndParse(commands Commands, parser Parser) (*Response, error) {
	resp, err := Fetch(commands.Commands()...)
	if err != nil {
		return nil, err
	}
	err = parser.Parse(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//FetchWithBodyAndParse fetch response and prase response with given commands ,body and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func FetchWithBodyAndParse(commands Commands, body io.Reader, parser Parser) (*Response, error) {
	preset := BuildPreset(commands.Commands()...)
	return FetchAndParse(preset.With(Body(body)), parser)
}

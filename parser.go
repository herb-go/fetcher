package fetcher

import (
	"encoding/json"
	"io"
)

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

//Should200 parser that check if status code == 200.
//Give parser will parse responese if success or respnse will be returned as error.
func Should200(p Parser) Parser {
	return ParserFunc(func(resp *Response) error {
		if resp.StatusCode != 200 {
			return resp
		}
		if p == nil {
			return nil
		}
		return p.Parse(resp)
	})
}

//ShouldSuccess parser that check if status code < 300.
//Give parser will parse responese if success or respnse will be returned as error.
func ShouldSuccess(p Parser) Parser {
	return ParserFunc(func(resp *Response) error {
		if resp.StatusCode >= 300 {
			return resp
		}
		if p == nil {
			return nil
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
		if p == nil {
			return nil
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

//FetchAndParse fetch request and prase response with given preset and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func FetchAndParse(preset *Preset, parser Parser) (*Response, error) {
	resp, err := Fetch(preset.Commands()...)
	if err != nil {
		return nil, err
	}
	err = parser.Parse(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//DoAndParse do request and prase response with given doer,preset and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func DoAndParse(doer Doer, preset *Preset, parser Parser) (*Response, error) {
	return FetchAndParse(preset.With(SetDoer(doer)), parser)
}

//FetchWithBodyAndParse fetch request and prase response with given preset ,body and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func FetchWithBodyAndParse(preset *Preset, body io.Reader, parser Parser) (*Response, error) {
	return FetchAndParse(preset.With(Body(body)), parser)
}

//FetchWithJSONBodyAndParse fetch request and prase response with given preset , body as json and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func FetchWithJSONBodyAndParse(preset *Preset, body interface{}, parser Parser) (*Response, error) {
	return FetchAndParse(preset.With(JSONBody(body)), parser)
}

//DoWithBodyAndParse do request and prase response with given doer,preset,body and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func DoWithBodyAndParse(doer Doer, preset *Preset, body io.Reader, parser Parser) (*Response, error) {
	return FetchAndParse(preset.With(SetDoer(doer), Body(body)), parser)
}

//DoWithJSONBodyAndParse do request and prase response with given doer,preset,body as json and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
func DoWithJSONBodyAndParse(doer Doer, preset *Preset, body interface{}, parser Parser) (*Response, error) {
	return FetchAndParse(preset.With(SetDoer(doer), JSONBody(body)), parser)
}

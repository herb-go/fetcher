package fetcher

import (
	"encoding/json"
	"io"
)

//Parser response parser interface
//Parser is responsible for closing response body.
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

//ShouldFunc create should parase func with given condition and next parser.
func ShouldFunc(condition func(resp *Response) error, next Parser) Parser {
	return ParserFunc(func(resp *Response) error {
		err := condition(resp)
		if err != nil && IsResponseErr(err) {
			resp.BodyContent()
			return err
		}
		if next == nil {
			return DefaultParser.Parse(resp)
		}
		return next.Parse(resp)
	})
}

//ShouldBetween parser that check if status code between min and max(include min/max).
//Give parser will parse responese if success or respnse will be returned as error.
func ShouldBetween(min int, max int, p Parser) Parser {
	return ShouldFunc(func(resp *Response) error {
		if resp.StatusCode < min || resp.StatusCode > max {
			return resp
		}
		return nil
	}, p)
}

//Should200 parser that check if status code == 200.
//Give parser will parse responese if success or respnse will be returned as error.
func Should200(p Parser) Parser {
	return ShouldFunc(func(resp *Response) error {
		if resp.StatusCode != 200 {
			return resp
		}
		return nil
	}, p)
}

//ShouldSuccess parser that check if status code < 300.
//Give parser will parse responese if success or respnse will be returned as error.
func ShouldSuccess(p Parser) Parser {
	return ShouldFunc(func(resp *Response) error {
		if resp.StatusCode >= 300 {
			return resp
		}
		return nil
	}, p)
}

//ShouldNoError parser that check if status code < 500.
//Give parser will parse responese if success or respnse will be returned as error.
func ShouldNoError(p Parser) Parser {
	return ShouldFunc(func(resp *Response) error {
		if resp.StatusCode >= 500 {
			return resp
		}
		return nil
	}, p)
}

//AsBytes create parser which parse givn byte slice from response.
func AsBytes(bytes *[]byte) Parser {
	return ParserFunc(func(resp *Response) error {
		bs, err := resp.BodyContent()
		if err != nil {
			return err
		}
		*bytes = make([]byte, len(bs))
		copy(*bytes, bs)
		return nil
	})
}

//AsDownload create parser which parse givn byte slice from response.
//You SHOULD NOT use BodyContent if you parsed response with Download Parser.
//This parser is designed to download file.
func AsDownload(w io.Writer) Parser {
	return ParserFunc(func(resp *Response) error {
		defer resp.Response.Body.Close()
		_, err := io.Copy(w, resp.Body)
		return err
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

var DefaultParser = AsBodyContent

//AsBodyContent read body to bodycontent
var AsBodyContent Parser = ParserFunc(func(resp *Response) error {
	_, err := resp.BodyContent()
	return err
})

//AsReader keep response body unreaded.
//You must close resp.Body manually.
var AsReader Parser = ParserFunc(func(resp *Response) error {
	return nil
})

//AsUselessBody do not use response body and close response body
var AsUselessBody Parser = ParserFunc(func(resp *Response) error {
	return resp.Body.Close()
})

//Fetch create new fetcher ,exec commands and fetch response.
//Return http response and any error if raised.
//Response returned will be parsed by defualt parser.
func Fetch(cmds ...Command) (*Response, error) {
	return FetchAndParse(Concat(cmds...), DefaultParser)
}

//FetchAndParse fetch request and prase response with given preset and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
//Response returned will be parsed by given parser or defualt parser if nill given.
func FetchAndParse(preset *Preset, parser Parser) (*Response, error) {
	resp, err := request(preset)
	if err != nil {
		return nil, err
	}
	err = parser.Parse(resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

//DoAndParse do request and prase response with given doer,preset and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
//Response returned will be parsed by given parser or defualt parser if nill given.
func DoAndParse(doer Doer, preset *Preset, parser Parser) (*Response, error) {
	return FetchAndParse(preset.Concat(SetDoer(doer)), parser)
}

//FetchWithBodyAndParse fetch request and prase response with given preset ,body and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
//Response returned will be parsed by given parser or defualt parser if nill given.
func FetchWithBodyAndParse(preset *Preset, body io.Reader, parser Parser) (*Response, error) {
	return FetchAndParse(preset.Concat(Body(body)), parser)
}

//FetchWithJSONBodyAndParse fetch request and prase response with given preset , body as json and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
//Response returned will be parsed by given parser or defualt parser if nill given.
func FetchWithJSONBodyAndParse(preset *Preset, body interface{}, parser Parser) (*Response, error) {
	return FetchAndParse(preset.Concat(JSONBody(body)), parser)
}

//DoWithBodyAndParse do request and prase response with given doer,preset,body and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
//Response returned will be parsed by given parser or defualt parser if nill given.
func DoWithBodyAndParse(doer Doer, preset *Preset, body io.Reader, parser Parser) (*Response, error) {
	return FetchAndParse(preset.Concat(SetDoer(doer), Body(body)), parser)
}

//DoWithJSONBodyAndParse do request and prase response with given doer,preset,body as json and parser if no error raised.
//Return response fetched and any error raised when fetching or parsing.
//Response returned will be parsed by given parser or defualt parser if nill given.
func DoWithJSONBodyAndParse(doer Doer, preset *Preset, body interface{}, parser Parser) (*Response, error) {
	return FetchAndParse(preset.Concat(SetDoer(doer), JSONBody(body)), parser)
}

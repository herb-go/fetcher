package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//ErrMsgLengthLimit max error message length
var ErrMsgLengthLimit = 512

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
	msg := fmt.Sprintf("fetcher:http error [%s %s ] %s : %s", r.Response.Request.Method, r.Response.Request.URL.String(), r.Status, string(bs))
	if len(msg) > ErrMsgLengthLimit {
		msg = msg[:ErrMsgLengthLimit]
	}
	return msg
}

//NewAPICodeErr make a api code error  which contains a error code.
func (r *Response) NewAPICodeErr(code interface{}) error {
	bs, err := r.BodyContent()
	if err != nil {
		return err
	}
	return NewAPICodeErr(r.Response.Request.URL.String(), r.Response.Request.Method, code, bs)

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

//NewAPICodeErr create a new api code error with given url,method,code,and content.
func NewAPICodeErr(url string, method string, code interface{}, content []byte) *APICodeErr {
	return &APICodeErr{
		URI:     url,
		Method:  method,
		Code:    fmt.Sprint(code),
		Content: content,
	}
}

//APICodeErr api code error struct.
type APICodeErr struct {
	//URI api uri.
	URI string
	//Code api error code.
	Code string
	//Method request method
	Method string
	//Content api response.
	Content []byte
}

//Error used as a error which return request url,request status,erro code,request content.
//Error max length is ErrMsgLengthLimit.
func (r *APICodeErr) Error() string {
	msg := fmt.Sprintf("fetcher:api error [%s %s] code %s : %s", r.URI, r.Method, r.Code, url.PathEscape(string(r.Content)[:ErrMsgLengthLimit]))
	if len(msg) > ErrMsgLengthLimit {
		msg = msg[:ErrMsgLengthLimit]
	}
	return msg
}

//GetAPIErrCode get api error code form error.
//Return empty string if err is not an ApiCodeErr
func GetAPIErrCode(err error) string {
	r, ok := err.(*APICodeErr)
	if ok {
		return r.Code
	}
	return ""
}

//GetAPIErrContent get api error code form error.
//Return empty string if err is not an ApiCodeErr
func GetAPIErrContent(err error) string {
	r, ok := err.(*APICodeErr)
	if ok {
		return string(r.Content)
	}
	return ""

}

//CompareAPIErrCode if check error is an ApiCodeErr with given api err code.
func CompareAPIErrCode(err error, code interface{}) bool {
	return GetAPIErrCode(err) == fmt.Sprint(code)
}

//IsResponseErr ehech if is response error
func IsResponseErr(err error) bool {
	_, ok := err.(*Response)
	return ok
}

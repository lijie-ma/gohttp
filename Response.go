package gohttp

import (
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
	Body string
	Error error
}

func (resp *Response) setBody() {
	b, e := ioutil.ReadAll(resp.Response.Body)
	defer resp.Response.Body.Close()
	if nil != e {
		resp.Error = e
		return
	}
	resp.Body = string(b)
}
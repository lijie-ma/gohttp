package httplib

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type HTTPLib struct {
	url      string
	method   string
	param    interface{}
	header   http.Header
	client   *http.Client
	request  *http.Request
	response *http.Response
}

func NewHTTPLib(url string, method string, param interface{}) *HTTPLib {
	return &HTTPLib{
		url:    url,
		method: strings.ToUpper(method),
		param:  param,
		client: &http.Client{
			Timeout: 0,
		},
	}
}

//DefaultCookieJar 设置默认jar
// domain 默认空
func DefaultCookieJar(domain ...string) *cookiejar.Jar {
	list := publicsuffix.List
	if 0 < len(domain) {
		list.PublicSuffix(domain[0])
	}
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: list})
	return jar
}

func (lib *HTTPLib) CookieJar() http.CookieJar {
	return lib.client.Jar
}

func (lib *HTTPLib) SetCookieJar(cookieJar *cookiejar.Jar) {
	lib.client.Jar = cookieJar
}

func (lib *HTTPLib) SetTimeout(timeout time.Duration) {
	lib.client.Timeout = timeout
}

func (lib *HTTPLib) SetProxy(proxy string) {

}

func (lib *HTTPLib) defaultHeader() http.Header {
	return http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
		"User-Agent":   []string{"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko"},
	}
}

func (lib *HTTPLib) CheckRedirect(fu func(req *http.Request, via []*http.Request) error) {
	lib.client.CheckRedirect = fu
}

func (lib *HTTPLib) SetHeader(header http.Header) {
	lib.header = header
}

func (lib *HTTPLib) Do() error {
	request, err := http.NewRequest(lib.method, lib.url, strings.NewReader(lib.getFormString()))
	if nil != err {
		return err
	}
	defer request.Body.Close()
	if 0 == len(lib.header) {
		request.Header = lib.defaultHeader()
	} else {
		request.Header = lib.header
	}
	lib.response, err = lib.client.Do(request)
	if nil != err {
		return err
	}
	return nil
}

func (lib *HTTPLib) Response() *http.Response {
	return lib.response
}

func (lib *HTTPLib) getFormString() string {
	if nil == lib.param {
		return ""
	}

	assertData, ok := lib.param.(map[string]interface{})
	if !ok {
		return lib.param.(string)
	}

	var urlValue = url.Values{}
	for k, v := range assertData {
		switch v.(type) {
		case float64:
			urlValue.Add(k, fmt.Sprintf("%f", v.(float64)))
		case int:
			urlValue.Add(k, strconv.Itoa(v.(int)))
		case int64:
			urlValue.Add(k, fmt.Sprintf("%d", v.(int64)))
		case string:
			urlValue.Add(k, v.(string))
		case []string:
			for _, sv := range v.([]string) {
				urlValue.Add(k+`[]`, sv)
			}
		case []int:
			for _, sv := range v.([]int) {
				urlValue.Add(k+`[]`, strconv.Itoa(sv))
			}
		}
	}
	return urlValue.Encode()
}

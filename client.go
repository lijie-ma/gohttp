package gohttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lijie-ma/utility"
	"golang.org/x/net/publicsuffix"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	client_version  = "0.0.1"
	COOKIES         = `cookies`
	HEADERS         = `headers`
	AUTH            = `auth`
	JSON            = `json`
	FORM_PARAMS     = `form_params`
	MULTIPART       = `multipart`
	MULTIPART_FILES = `files`
	BASE_URI        = `base_uri`
	PROXY           = `proxy`
	TIMEOUT         = `timeout`
	QUERY           = `query`
)

var (
	errTypetimeout = errors.New("invalid timeout type, require int or time.Duration")
	errTypeQuery   = errors.New("invalid query type, require string")
	errEmptyURI    = errors.New("empty base_uri set")
	errTypeURI     = errors.New("invalid base_uri type, require string")
)

type Client struct {
	config            map[string]interface{}
	uri               *url.URL
	p                 *sync.Pool
	httpClient        *http.Client
	currentHttpClient *http.Client
	errs              []error
}

func NewClient(config map[string]interface{}) *Client {
	c := &Client{errs: make([]error, 0, 2)}
	c.configDefault(config)
	c.setHttpClient()
	c.p = &sync.Pool{New: func() interface{} {
		return c.httpClient
	}}
	return c
}

func (c *Client) configDefault(config map[string]interface{}) {
	c.defaultHeaders(config)
	if _, ok := config[COOKIES]; !ok {
		config[COOKIES] = false
	}
	c.config = config
}

func (c *Client) defaultHeaders(config map[string]interface{}) {
	headers := http.Header{}
	if _, ok := config[HEADERS]; !ok {
		headers.Add("User-Agent", defaultUserAgent())
	} else {
		tmp := config[HEADERS]
		switch tmp.(type) {
		case []map[string]string:
			for _, hh := range tmp.([]map[string]string) {
				for k, v := range hh {
					headers.Add(k, v)
				}
			}
		case map[string][]string:
			headers = http.Header(tmp.(map[string][]string))
		case http.Header:
			headers = tmp.(http.Header)
		}
		if 0 == len(headers.Get("User-Agent")) && 0 == len(headers.Get("user-agent")) {
			headers.Add("User-Agent", defaultUserAgent())
		}
	}
	config[HEADERS] = headers
}

func (c *Client) setHttpClient() {
	c.httpClient = &http.Client{Timeout: 0 * time.Second}
	if v, ok := c.config[TIMEOUT]; ok {
		switch v.(type) {
		case int:
			c.httpClient.Timeout = time.Duration(v.(int)) * time.Second
		case time.Duration:
			c.httpClient.Timeout = v.(time.Duration)
		default:
			c.addError(errTypetimeout)
		}
	}
	if rawurl, ok := c.config[PROXY]; ok {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(rawurl.(string))
		}
		c.httpClient.Transport = &http.Transport{Proxy: proxy}
	}
}

//option中的设置会覆盖全局的client中的设置
// 并返回新设置client
func (c *Client) getHttpClient(options map[string]interface{}) *http.Client {
	c.currentHttpClient = c.p.Get().(*http.Client)
	if v, ok := options[TIMEOUT]; ok {
		switch v.(type) {
		case int:
			c.currentHttpClient.Timeout = time.Duration(v.(int)) * time.Second
		case time.Duration:
			c.currentHttpClient.Timeout = v.(time.Duration)
		default:
			c.addError(errTypetimeout)
		}
	}
	if rawurl, ok := options[PROXY]; ok {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(rawurl.(string))
		}
		c.currentHttpClient.Transport = &http.Transport{Proxy: proxy}
	}

	return c.currentHttpClient
}

// ResetErrors每次请求前重置 可以通过配置 reset_error 的值为 false 来禁止
func (c *Client) ResetErrors() {
	c.resetErrors(nil)
}

func (c *Client) resetErrors(option map[string]interface{}) {
	tempFunc := func(v interface{}) bool {
		switch v.(type) {
		case bool:
			return v.(bool)
		}
		return false
	}
	tmpkey := "reset_error"
	if v, ok := option[tmpkey]; ok {
		if tempFunc(v) {
			c.errs = c.errs[:0]
		}
		return
	}
	if v, ok := c.config[tmpkey]; ok {
		if tempFunc(v) {
			c.errs = c.errs[:0]
		}
		return
	}
	c.errs = c.errs[:0]
}

func (c *Client) GetErrors() []error {
	return c.errs
}

func (c *Client) addError(e error) {
	c.errs = append(c.errs, e)
}

func (c *Client) request(method string, uri string, options map[string]interface{}) *Response {
	c.resetErrors(options)
	request, err := http.NewRequest(method, c.rebuildURI(uri, options), c.requestBody(options))
	if nil != err {
		c.addError(err)
		return nil
	}

	c.setRequestHeader(request, options)
	if 0 < len(c.errs) {
		return nil
	}
	httpClient := c.getHttpClient(options)
	c.setCookies(httpClient, options)
	response, err := httpClient.Do(request)
	c.p.Put(httpClient)
	if nil != err {
		c.addError(err)
		return nil
	}
	resp := &Response{}
	resp.Response = response
	resp.setBody()
	return resp
}

func (c *Client) GetCookies() []*http.Cookie {
	if nil != c.currentHttpClient.Jar {
		return c.currentHttpClient.Jar.Cookies(c.uri)
	}
	return nil
}

func (c *Client) CloseCookies() {
	header, ok := c.config[HEADERS].(http.Header)
	if ok {
		header.Del("Cookie")
	}
	c.config[COOKIES] = false
}

func (c *Client) setCookies(client *http.Client, options map[string]interface{}) {
	setFunc := func(client *http.Client, cookies interface{}) {
		switch cookies.(type) {
		case bool:
			if cookies.(bool) {
				client.Jar = DefaultCookieJar(c.uri.Host)
			} else if nil != client.Jar {
				client.Jar = nil
			}
		case []*http.Cookie:
			if 0 == len(cookies.([]*http.Cookie)) {
				break
			}
			client.Jar = DefaultCookieJar(c.uri.Host)
			client.Jar.SetCookies(c.uri, cookies.([]*http.Cookie))
		}
	}

	if v, ok := options[COOKIES]; ok {
		if nil != v {
			setFunc(client, v)
		}

	} else if v, ok := c.config[COOKIES]; ok {
		if nil != v {
			setFunc(client, v)
		}
	}
}

func (c *Client) rebuildURI(uri string, option map[string]interface{}) string {
	queryStr := ``
	if v, ok := option[QUERY]; ok {
		switch v.(type) {
		case string:
			queryStr = v.(string)
		default:
			c.addError(errTypeQuery)
		}
	} else if v, ok := c.config[QUERY]; ok {
		switch v.(type) {
		case string:
			queryStr = v.(string)
		default:
			c.addError(errTypeQuery)
		}
	}
	if 0 == strings.Index(uri, `http://`) || 0 == strings.Index(uri, `https://`) {
		uriParse, err := url.Parse(uri)
		if nil != err {
			c.addError(err)
			return ``
		}
		c.uri = uriParse
		if 0 == len(queryStr) {
			return uri
		}
		if -1 == strings.Index(uri, `?`) {
			uri += `?`
		} else {
			uri += `&`
		}
		return uri + queryStr
	} else if 0 < len(queryStr) {
		if -1 == strings.Index(uri, `?`) {
			uri += `?`
		} else {
			uri += `&`
		}
		uri += queryStr
	}
	baseUri, ok := c.config[BASE_URI]
	if !ok {
		c.addError(errEmptyURI)
		return ``
	}

	switch baseUri.(type) {
	case string:
	default:
		c.addError(errTypeURI)
		return ``
	}
	uriParse, err := url.Parse(baseUri.(string))
	if nil != err {
		c.addError(err)
		return ``
	}
	c.uri = uriParse
	return c.uri.Scheme + `://` + c.uri.Host + `/` + strings.TrimLeft(uri, `/`)
}

func (c *Client) setRequestHeader(r *http.Request, option map[string]interface{}) {
	if h, ok := c.config[HEADERS]; ok {
		r.Header = h.(http.Header)
	}
	if v, ok := c.config[AUTH]; ok {
		r.SetBasicAuth(v.([]string)[0], v.([]string)[1])
	}
}

func (c *Client) requestBody(option map[string]interface{}) io.Reader {
	funcJson := func(v interface{}) io.Reader {
		h := c.config["headers"].(http.Header)
		h.Set("Content-Type", "application/json")
		s := ""
		switch v.(type) {
		case string:
			s = v.(string)
		default:
			b, e := json.Marshal(v)
			if nil != e {
				c.addError(e)
			}
			s = string(b)
		}
		return strings.NewReader(s)
	}
	if v, ok := option["json"]; ok {
		return funcJson(v)
	}
	if v, ok := option["form_params"]; ok {
		h := c.config["headers"].(http.Header)
		h.Set("Content-Type", "application/x-www-form-urlencoded")
		return strings.NewReader(utility.HttpBuildQuery(v.(map[string]interface{})))
	}
	if v, ok := option[MULTIPART]; ok {
		return c.setUploads(v.(map[string]interface{}))
	}
	if v, ok := c.config["json"]; ok {
		return funcJson(v)
	}
	if v, ok := c.config["form_params"]; ok {
		h := c.config["headers"].(http.Header)
		h.Set("Content-Type", "application/x-www-form-urlencoded")
		return strings.NewReader(utility.HttpBuildQuery(v.(map[string]interface{})))
	}
	if v, ok := c.config[MULTIPART]; ok {
		return c.setUploads(v.(map[string]interface{}))
	}
	return nil
}

// 设置上传信息
func (c *Client) setUploads(uploads map[string]interface{}) io.Reader {
	h := c.config[HEADERS].(http.Header)
	h.Set("Content-Type", "multipart/form-data")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	files, ok := uploads[MULTIPART_FILES]
	if ok {
		for field, file := range files.(map[string]string) {
			fp, err := os.Open(file)
			if err != nil {
				c.addError(err)
				continue
			}
			defer fp.Close()
			part, err := writer.CreateFormFile(field, filepath.Base(file))
			if nil != err {
				c.addError(err)
				continue
			}
			_, err = io.Copy(part, fp)
			if nil != err {
				c.addError(err)
				continue
			}
			h.Set("Content-Type", writer.FormDataContentType())
		}
	}
	fileds, ok := uploads[FORM_PARAMS]
	if ok {
		for field, value := range fileds.(map[string]string) {
			writer.WriteField(field, value)
		}
	}

	err := writer.Close()
	if nil != err {
		c.addError(err)
		return nil
	}

	return body
}

// 发送post 请求，如果出错，则返回nil， 可以通过 GetErrors 拿到错误信息
// options 可以设置 表单信息
//  			form_params
//				uploads
//				json
// 这些会覆盖全局的
func (c *Client) Post(uri string, options map[string]interface{}) *Response {
	return c.request("POST", uri, options)
}

//Get 如果出错，则返回nil， 可以通过 GetErrors 拿到错误信息
func (c *Client) Get(uri string, options map[string]interface{}) *Response {
	return c.request("GET", uri, options)
}

func (c *Client) Head(uri string, options map[string]interface{}) *Response {
	return c.request("HEAD", uri, options)
}

func Post(uri string, options map[string]interface{}) *Response {
	return NewClient(map[string]interface{}{}).Post(uri, options)
}

func Get(uri string, options map[string]interface{}) *Response {
	return NewClient(map[string]interface{}{}).Get(uri, options)
}

func Head(uri string, options map[string]interface{}) *Response {
	return NewClient(map[string]interface{}{}).Head(uri, options)
}

func defaultUserAgent() string {
	return "gohttp/" + client_version + "  golang " + runtime.Version()
}

func DefaultCookieJar(domain ...string) *cookiejar.Jar {
	list := publicsuffix.List
	if 0 < len(domain) {
		list.PublicSuffix(domain[0])
	}
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: list})
	return jar
}

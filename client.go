package gohttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lijie-ma/utility"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	client_version = "0.0.1"
)

var (
	errTypetimeout = errors.New("invalid timeout type, require int or time.Duration")
	errTypeQuery   = errors.New("invalid query type, require string")
	errEmptyURI    = errors.New("empty base_uri set")
	errTypeURI     = errors.New("invalid base_uri type, require string")
)

type Client struct {
	config map[string]interface{}
	errs   []error
}

func NewClient(config map[string]interface{}) *Client {
	c := &Client{errs: make([]error, 0, 2)}
	c.configDefault(config)
	return c
}

func (c *Client) configDefault(config map[string]interface{}) {
	c.defaultHeaders(config)
	if _, ok := config["cookies"]; !ok {
		config["cookies"] = false
	}
	if _, ok := config["http_errors"]; !ok {
		config["http_errors"] = true
	}
	c.config = config
}

func (c *Client) defaultHeaders(config map[string]interface{}) {
	headers := http.Header{}
	if _, ok := config["headers"]; !ok {
		headers.Add("User-Agent", defaultUserAgent())
		config["headers"] = headers
	} else {
		tmp := config["headers"]
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
		config["headers"] = headers
	}
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

func (c *Client) request(method string, uri string, option map[string]interface{}) *Response {
	c.resetErrors(option)
	request, err := http.NewRequest(method, c.rebuildURI(uri, option), c.requestBody(option))
	if nil != err {
		c.addError(err)
		return nil
	}
	c.setRequestHeader(request, option)
	if 0 < len(c.errs) {
		return nil
	}
	response, err := c.getHttpClient(option).Do(request)
	if nil != err {
		c.addError(err)
		return nil
	}
	resp := &Response{}
	resp.Response = response
	resp.setBody()
	return resp
}

func (c *Client) getHttpClient(option map[string]interface{}) *http.Client {
	clientHttp := &http.Client{Timeout: 0 * time.Second}
	if v, ok := c.config["timeout"]; ok {
		switch v.(type) {
		case int:
			clientHttp.Timeout = time.Duration(v.(int)) * time.Second
		case time.Duration:
			clientHttp.Timeout = v.(time.Duration)
		default:
			c.addError(errTypetimeout)
		}
	}
	if rawurl, ok := c.config["proxy"]; ok {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(rawurl.(string))
		}
		clientHttp.Transport = &http.Transport{Proxy: proxy}
	}

	return clientHttp
}

func (c *Client) rebuildURI(uri string, option map[string]interface{}) string {
	queryStr := ``
	if v, ok := option["query"]; ok {
		switch v.(type) {
		case string:
			queryStr = v.(string)
		default:
			c.addError(errTypeQuery)
		}
	}
	if 0 == strings.Index(uri, `http://`) || 0 == strings.Index(uri, `https://`) {
		if 0 < len(queryStr) {
			return uri
		}
		if -1 == strings.Index(uri, `?`) {
			uri += `?`
		}
		return uri + queryStr
	} else if 0 < len(queryStr) {
		if -1 == strings.Index(uri, `?`) {
			uri += `?`
		}
		uri += queryStr
	}
	baseUri, ok := c.config["base_uri"]
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
	return uriParse.Scheme + `://` + uriParse.Host + `/` + strings.TrimLeft(uri, `/`)
}

func (c *Client) setRequestHeader(r *http.Request, option map[string]interface{}) {
	if v, ok := c.config["auth"]; ok {
		r.SetBasicAuth(v.([]string)[0], v.([]string)[1])
	}
	if h, ok := c.config["headers"]; ok {
		r.Header = h.(http.Header)
	}
}

func (c *Client) requestBody(option map[string]interface{}) io.Reader {
	if v, ok := c.config["json"]; ok {
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
	if v, ok := c.config["form_params"]; ok {
		h := c.config["headers"].(http.Header)
		h.Set("Content-Type", "application/x-www-form-urlencoded")
		return strings.NewReader(utility.HttpBuildQuery(v.(map[string]interface{})))
	}

	return c.setUploads(option)
}

// 设置上传信息
func (c *Client) setUploads(option map[string]interface{}) io.Reader {
	uploads, ok := c.config["uploads"]
	if !ok {
		return nil
	}
	h := c.config["headers"].(http.Header)
	h.Set("Content-Type", "multipart/form-data")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	files, ok := uploads.(map[string]interface{})["files"]
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
	fileds, ok := uploads.(map[string]interface{})["form_params"]
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

func (c *Client) prepareDefaults(option map[string]interface{}) map[string]interface{} {
	return option
}

// 发送post 请求，如果出错，则返回nil， 可以通过 GetErrors 拿到错误信息
func (c *Client) Post(uri string, options map[string]interface{}) *Response {
	return c.request("POST", uri, options)
}

//Get 如果出错，则返回nil， 可以通过 GetErrors 拿到错误信息
func (c *Client) Get(uri string, options map[string]interface{}) *Response {
	return c.request("GET", uri, options)
}

func (c *Client) Put() {

}

func (c *Client) Delete() {

}
func (c *Client) Head() {

}

func defaultUserAgent() string {
	return "gohttp/" + client_version + "  golang " + runtime.Version()
}

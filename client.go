package gohttp

import (
	"bytes"
	"encoding/json"
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

type Client struct {
	config map[string]interface{}
}

func NewClient(config map[string]interface{}) *Client {
	c := &Client{}
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

func (c *Client) request(method string, uri string, option map[string]interface{}) (*Response, error) {
	request, err := http.NewRequest(method, uri, c.requestBody(option))
	if nil != err {
		return nil, err
	}
	c.setRequestHeader(request, option)
	response, err := c.getHttpClient(option).Do(request)
	if nil != err {
		return nil, err
	}
	resp := &Response{}
	resp.Response = response
	resp.setBody()
	return resp, err
}

func (c *Client) getHttpClient(option map[string]interface{}) *http.Client {
	clientHttp := &http.Client{Timeout: 0 * time.Second}
	if v, ok := c.config["timeout"]; ok {
		switch v.(type) {
		case int:
			clientHttp.Timeout = time.Duration(v.(int)) * time.Second
		case time.Duration:
			clientHttp.Timeout = v.(time.Duration)
		case float32:
			clientHttp.Timeout = time.Duration(v.(float32)) * time.Second
		case float64:
			clientHttp.Timeout = time.Duration(v.(float64)) * time.Second
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
				panic(e)
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
				panic(err)
			}
			defer fp.Close()
			part, err := writer.CreateFormFile(field, filepath.Base(file))
			if nil != err {
				panic(err)
			}
			_, err = io.Copy(part, fp)
			if nil != err {
				panic(err)
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
		return nil
	}

	return body
}

func (c *Client) prepareDefaults(option map[string]interface{}) map[string]interface{} {
	return option
}

func (c *Client) Post(uri string, options map[string]interface{}) (*Response, error) {
	return c.request("POST", uri, options)
}

func (c *Client) Get(uri string, options map[string]interface{}) (*Response, error) {
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

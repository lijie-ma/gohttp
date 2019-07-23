package gohttp

import (
	"io/ioutil"
	"net/http"
	"testing"
)


func TestClient_Post(t *testing.T) {
	v := map[string]interface{}{
		"form_params" : map[string]interface{}{
			"key" :"ivideo_index",
		},
		"proxy": "http://127.0.0.1:8888",
	}
	c := NewClient(v)
	resp, err := c.Post("http://10.16.155.5:8090/cms/getone",nil)
	t.Log(err)
	b, e := ioutil.ReadAll(resp.Body)
	t.Log(string(b), e)
}

func TestGet(t *testing.T) {
	config := map[string]interface{}{
		"headers" : map[string][]string {
			"ua": []string{"cc", "ddd"},
		},
	}
	headers := http.Header{}
	if _, ok := config["headers"]; !ok {
		headers.Set("User-Agent", defaultUserAgent())
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
		}
		delete(config, "headers")
		config["headers"] = headers
	}
	t.Log(config)
}

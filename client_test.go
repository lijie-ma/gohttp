package gohttp

import (
	"net/http"
	"testing"
)

func TestClient_Post(t *testing.T) {
	v := map[string]interface{}{
		"form_params": map[string]interface{}{
			"key": "ivideo_index",
		},
		"proxy": "http://127.0.0.1:8888",
	}
	c := NewClient(v)
	resp, err := c.Post("http://10.16.155.5:8090/cms/getone", nil)
	t.Log(err)
	t.Log(resp.Body)

	v2 := map[string]interface{}{
		"json": map[string]interface{}{
			"key": "ivideo_index",
		},
		"proxy": "http://127.0.0.1:8888",
	}
	c2 := NewClient(v2)
	resp, err = c2.Post("http://192.168.56.102/test.php", nil)
	t.Log(err)
	t.Log(resp.Body)

	v3 := map[string]interface{}{
		"uploads": map[string]interface{}{
			"files" : map[string]string {
				"f1": "D:/download/2.sql",
				"f2": "D:/download/1.sql",
			},
			"form_params" : map[string]string {
				"key": "ivideo_index",
			},
		},
		"proxy": "http://127.0.0.1:8888",
	}
	c3 := NewClient(v3)
	resp, err = c3.Post("http://192.168.56.102/upload.php", nil)
	t.Log(err)
	t.Log(resp.Body)
}

func TestGet(t *testing.T) {
	config := map[string]interface{}{
		"headers": map[string][]string{
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

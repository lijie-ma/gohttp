package gohttp

import (
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
	resp := c.Post("http://10.16.155.5:8090/cms/getone", nil)
	t.Log(resp.Body)

	v2 := map[string]interface{}{
		"json": map[string]interface{}{
			"key": "ivideo_index",
		},
		"proxy": "http://127.0.0.1:8888",
	}
	c2 := NewClient(v2)
	resp = c2.Post("http://192.168.56.102/test.php", nil)
	t.Log(resp.Body)

	v3 := map[string]interface{}{
		"uploads": map[string]interface{}{
			"files": map[string]string{
				"f1": "D:/download/2.sql",
				"f2": "D:/download/1.sql",
			},
			"form_params": map[string]string{
				"key": "ivideo_index",
			},
		},
		"proxy": "http://127.0.0.1:8888",
	}
	c3 := NewClient(v3)
	resp = c3.Post("http://192.168.56.102/upload.php", nil)
	t.Log(resp.Body)
}

func TestBaseUri(t *testing.T) {
	v2 := map[string]interface{}{
		"json": map[string]interface{}{
			"key": "ivideo_index",
		},
		"proxy":    "http://127.0.0.1:8888",
		"base_uri": "http://192.168.56.102/",
	}
	c2 := NewClient(v2)
	resp := c2.Post("/test.php", nil)
	t.Log("base_uri\t", resp.Body)
}

func TestOption(t *testing.T) {
	v := map[string]interface{}{
		"proxy": "http://127.0.0.1:8888",
	}
	option := map[string]interface{}{
		"form_params": map[string]interface{}{
			"key": "ivideo_index",
		},
	}
	c := NewClient(v)
	resp := c.Post("http://10.16.155.5:8090/cms/getone", option)
	t.Log("option\t", resp.Body)

	v2 := map[string]interface{}{
		"json": map[string]interface{}{
			"key": "ivideo_index",
		},
		"proxy": "http://127.0.0.1:8888",
	}
	option2 := map[string]interface{}{
		"json": `{"key":"value"}`,
	}
	c2 := NewClient(v2)
	resp = c2.Post("http://192.168.56.102/test.php", option2)
	t.Log("option\t", resp.Body)

	v3 := map[string]interface{}{
		"uploads": map[string]interface{}{
			"files": map[string]string{
				"f1": "D:/download/2.sql",
				"f2": "D:/download/1.sql",
			},
			"form_params": map[string]string{
				"key": "ivideo_index",
			},
		},
		"proxy": "http://127.0.0.1:8888",
	}
	option3 := map[string]interface{}{
		"uploads": map[string]interface{}{
			"files": map[string]string{
				"f1": "D:/download/2.sql",
			},
			"form_params": map[string]string{
				"key": "ivideo_index",
			},
		},
		"proxy": "127.0.0.1:8888",
	}
	c3 := NewClient(v3)
	resp = c3.Post("http://192.168.56.102/upload.php", option3)
	t.Log("option\t", resp.Body)
}

func TestCookie(t *testing.T) {
	v2 := map[string]interface{}{
		"json": map[string]interface{}{
			"key": "ivideo_index",
		},
		"proxy": "http://127.0.0.1:8888",
		"cookies": true,
	}
	option2 := map[string]interface{}{
		"json": `{"key":"value"}`,
	}
	c2 := NewClient(v2)
	resp := c2.Post("http://192.168.56.102/s1.php", nil)

	option2["cookies"] = resp.Cookies()
	resp = c2.Post("http://192.168.56.102/s2.php", option2)
	t.Log("Cookie\t", resp.Body)
	////close cookies
	//option2["cookies"] = false
	//c2.CloseCookies()
	//resp = c2.Post("http://192.168.56.102/s2.php", option2)
	//t.Log("Cookie colse\t", resp.Body)
}

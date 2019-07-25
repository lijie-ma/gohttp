# gohttp
go http client
仿照php中GuzzleHttp类库翻写
## 安装
```
go get -u github.com/lijie-ma/gohttp
```

## demo

### post请求
```golang

v := map[string]interface{}{
    "form_params": map[string]interface{}{
        "key": "ivideo_index",
    },
    "proxy": "http://127.0.0.1:8888",
}
c := NewClient(v)
resp, err := c.Post("http://10.16.155.5:8090/cms/getone", nil)
log.Println(resp, err)

```
### 文件上传
```cassandraql
v3 := map[string]interface{}{
    "uploads": map[string]interface{}{
        "files" : map[string]string {         //注意类型map[string]string
            "f1": "D:/download/2.sql",
            "f2": "D:/download/1.sql",
        },
        "form_params" : map[string]string {  //注意类型map[string]string
            "key": "ivideo_index",
        },
    },
    "proxy": "http://127.0.0.1:8888",
}
c3 := NewClient(v3)
resp, err = c3.Post("http://192.168.56.102/upload.php", nil)
fmt.Println(resp.Body)
```

### login (代码见test)
```cassandraql
v := map[string]interface{}{
    "form_params": map[string]interface{}{
        "name": "aa",
    },
    "proxy":   proxy,
    "cookies": true,
}
option2 := map[string]interface{}{
    "json": `{"key":"value"}`,
}
c2 := NewClient(v)
resp := c2.Post(loginUri, nil)
option2["cookies"] = c2.GetCookies()
fmt.Println(resp.Body, resp.Cookies(), c2.GetCookies())
resp = c2.Post(session2Uri, option2)
t.Log("Cookie\t", resp.Body)

```

## doc
### 请求选项
| 选项 | 类型 | 样例 | 备注 |
| :------:| :------: | :------: | :------: |
| auth | []string | [name, password] | |
| cookies | bool &vert; []*http.Cookies | map[string]interface{}{"cookies":true} or <br>map[string]interface{}{"cookies":[]*http.Cookies{}}| true 表示开启cookie， http.cookies 表示本次请求要传送的cookies|
| connect_timeout | int &vert; time.Duration |  | 连接超时时间，如果传递int类型，单位是秒 |
| form_params | map[string]interface{} |  |关联数组由表单字段键值对构成，每个字段值可以是一个字符串或一个包含字符串元素的数组。 当没有准备 "Content-Type" 报文头的时候，将设置为 "application/x-www-form-urlencoded" |
| headers |[]map[string]string &vert; http.Header |  |默认的设置有ua、content-type |
| proxy | string |  | reqeust时 使用的代理 |
| query | string |  | 会被拼接到url后，目前仅支持string类型 |
| json | string &vert; interface{} |  | interface{}类型会调用json.Marshal， content-type 会设置为application/json |
| uploads | map[string]interface{} |  | 参考类domo 中的uploads[###文件上传](文件上传)  |
# gohttp
go http client
仿照php中GuzzleHttp类库翻写
## 安装
```
go get -u github.com/lijie-ma/gohttp
```

## demo

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

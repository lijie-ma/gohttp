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
###文件上传
```cassandraql
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
fmt.Println(resp.Body)
```

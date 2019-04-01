# httplib
go 模拟登录

## 安装
```
go get -u github.com/lijie-ma/httplib
```

## demo

```golang
jar := httplib.DefaultCookieJar()
httpClient := httplib.NewHTTPLib(`http://127.0.0.1/login`, http.MethodPost,
		map[string]interface{}{"id": id, "password": pwd})
httpClient.SetCookieJar(jar)
err := httpClient.Do()
if nil != err {
	return
}

```

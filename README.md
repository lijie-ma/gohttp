# go-httplib
go 模拟登录

## demo

```golang
jar := httplib.DefaultCookieJar()
httpClient := httplib.NewHTTPLib(ads.spider.Login, http.MethodPost,
		map[string]interface{}{"id": id, "password": pwd})
httpClient.SetCookieJar(jar)
httpClient.SetHeader(http.Header{
	"Content-Type": []string{"application/x-www-form-urlencoded"},
	"User-Agent":   []string{"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko"},
})
err := httpClient.Do()
if nil != err {
	return
}

```

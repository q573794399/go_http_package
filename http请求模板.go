package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Client = Init("") //在这里设置代理 不设置就 传入空字符串  设置就传入代理地址

func Init(proxy_ip string) *http.Client {
	//初始化
	u, _ := url.Parse(proxy_ip)
	var t *http.Transport
	if proxy_ip == "" {
		t = &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 10,
			IdleConnTimeout: time.Duration(10) * time.Second,
		}
	} else {
		t = &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 10,
			IdleConnTimeout: time.Duration(10) * time.Second,
			Proxy:           http.ProxyURL(u),
		}
	}
	client := &http.Client{
		Transport: t,
		Timeout:   time.Duration(5) * time.Second,
	} //创建一个连接

	return client
}

//普通get请求 返回字符串  用于爬取网页
func Get_Text(link_url string, head map[string]string) (string, error) {
	//普通get请求
	req, _ := http.NewRequest("GET", link_url, nil) //配置本次请求的方式和url

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return "", errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return "", errors.New("状态码异常")
	}

	defer resp.Body.Close()                 //关闭连接
	body, err2 := ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err2 != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err2)
		return "", errors.New("数据转换为字符串错误 一般是返回的数据是压缩的 需要在请求头里面设置接口明码返回")
		//	设置接口明码返回  head["Accept-Encoding"] = "identity"
	}

	return string(body), nil //将数据转换为字符串类型 并且作为函数的返回值
}

//普通get请求 返回数据流 用于下载文件
func Get_Byte(link_url string, head map[string]string) ([]byte, error) {
	//普通get请求 返回数据流 用于下载文件
	req, _ := http.NewRequest("GET", link_url, nil) //配置本次请求的方式和url
	var body []byte

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return body, errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return body, errors.New("状态码异常")
	}
	defer resp.Body.Close()               //关闭连接
	body, err = ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err != nil {
		fmt.Printf("err ByteGet2: %v\n", err)
		return nil, errors.New("数据转换错误")
	}
	return body, nil // 并且作为函数的返回值
}

//普通get请求 返回json数据
func Get_Json(link_url string, head map[string]string) (map[string]interface{}, error) {
	//普通get请求 返回json数据
	req, _ := http.NewRequest("GET", link_url, nil) //配置本次请求的方式和url
	var r map[string]interface{}
	r = make(map[string]interface{})

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return r, errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return r, errors.New("状态码异常")
	}

	defer resp.Body.Close()                //关闭连接
	body, err := ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err)
		return r, errors.New("数据转换为字符串错误 一般是返回的数据是压缩的 需要在请求头里面设置接口明码返回")
		//	设置接口明码返回  head["Accept-Encoding"] = "identity"
	}

	r = make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err)
		return r, errors.New("数据转换为json格式失败")
	}

	return r, nil //将数据转换为字符串类型 并且作为函数的返回值
}

//普通的post 请求 返回网页
func Post_From_Text(link_url string, head, data map[string]string) (string, error) {

	urlValues := url.Values{}    //创建提交参数表单
	for key, val := range data { //遍历map
		urlValues.Add(key, val) //将内容写入到表单内容
	}

	req, _ := http.NewRequest("POST", link_url, strings.NewReader(urlValues.Encode())) //配置本次请求的方式和url

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //通过设置请求头参数 让数据变成表单form类型进行提交  默认的提交的是data

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return "", errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return "", errors.New("状态码异常")
	}

	defer resp.Body.Close()                 //关闭连接
	body, err2 := ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err2 != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err2)
		return "", errors.New("数据转换为字符串错误 一般是返回的数据是压缩的 需要在请求头里面设置接口明码返回")
		//	设置接口明码返回  head["Accept-Encoding"] = "identity"
	}

	return string(body), nil //将数据转换为字符串类型 并且作为函数的返回值

}

//普通的post请求 返回json
func Post_From_Json(link_url string, head, data map[string]string) (map[string]interface{}, error) {

	var r map[string]interface{}
	r = make(map[string]interface{})

	urlValues := url.Values{}    //创建提交参数表单
	for key, val := range data { //遍历map
		urlValues.Add(key, val) //将内容写入到表单内容
	}

	req, _ := http.NewRequest("POST", link_url, strings.NewReader(urlValues.Encode())) //配置本次请求的方式和url

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //通过设置请求头参数 让数据变成表单form类型进行提交  默认的提交的是data

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return r, errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return r, errors.New("状态码异常")
	}

	defer resp.Body.Close()                 //关闭连接
	body, err2 := ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err2 != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err2)
		return r, errors.New("数据转换为字符串错误 一般是返回的数据是压缩的 需要在请求头里面设置接口明码返回")
		//	设置接口明码返回  head["Accept-Encoding"] = "identity"
	}

	r = make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err)
		return r, errors.New("数据转换为json格式失败")
	}

	return r, nil //将数据转换为字符串类型 并且作为函数的返回值

}

//普通的post请求 返回数据流 用于下载文件
func Post_From_Byte(link_url string, head, data map[string]string) ([]byte, error) {
	var body []byte

	urlValues := url.Values{}    //创建提交参数表单
	for key, val := range data { //遍历map
		urlValues.Add(key, val) //将内容写入到表单内容
	}

	req, _ := http.NewRequest("POST", link_url, strings.NewReader(urlValues.Encode())) //配置本次请求的方式和url

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //通过设置请求头参数 让数据变成表单form类型进行提交  默认的提交的是data

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return body, errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return body, errors.New("状态码异常")
	}

	defer resp.Body.Close()                 //关闭连接
	body, err2 := ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err2 != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err2)
		return body, errors.New("数据转换为字符串错误 一般是返回的数据是压缩的 需要在请求头里面设置接口明码返回")
		//	设置接口明码返回  head["Accept-Encoding"] = "identity"
	}

	return body, nil //将数据转换为字符串类型 并且作为函数的返回值

}

//json的post请求  返回网页
func Post_Json_Text(link_url string, head, data map[string]string) (string, error) {

	respdata, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return "", errors.New("提交参数转换为json失败")
	}

	req, _ := http.NewRequest("POST", link_url, bytes.NewReader(respdata)) //配置本次请求的方式和url

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return "", errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return "", errors.New("状态码异常")
	}

	defer resp.Body.Close()                 //关闭连接
	body, err2 := ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err2 != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err2)
		return "", errors.New("数据转换为字符串错误 一般是返回的数据是压缩的 需要在请求头里面设置接口明码返回")
		//	设置接口明码返回  head["Accept-Encoding"] = "identity"
	}

	return string(body), nil //将数据转换为字符串类型 并且作为函数的返回值

}

//json的post请求  返回json
func Post_Json_Json(link_url string, head, data map[string]string) (map[string]interface{}, error) {

	var r map[string]interface{}
	r = make(map[string]interface{})

	respdata, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return r, errors.New("提交参数转换为json失败")
	}

	req, _ := http.NewRequest("POST", link_url, bytes.NewReader(respdata)) //配置本次请求的方式和url

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return r, errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return r, errors.New("状态码异常")
	}

	defer resp.Body.Close()                 //关闭连接
	body, err2 := ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err2 != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err2)
		return r, errors.New("数据转换为字符串错误 一般是返回的数据是压缩的 需要在请求头里面设置接口明码返回")
		//	设置接口明码返回  head["Accept-Encoding"] = "identity"
	}

	r = make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err)
		return r, errors.New("数据转换为json格式失败")
	}

	return r, nil //将数据转换为字符串类型 并且作为函数的返回值

}

//json的post请求 返回数据流 用于下载文件
func Post_Json_Byte(link_url string, head, data map[string]string) ([]byte, error) {
	var body []byte

	respdata, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return body, errors.New("提交参数转换为json失败")
	}

	req, _ := http.NewRequest("POST", link_url, bytes.NewReader(respdata)) //配置本次请求的方式和url

	//设置请求头内容
	for key, val := range head { //遍历map
		req.Header.Set(key, val)
	}

	resp, err := Client.Do(req) //设置完毕后 发出网络请求
	if err != nil {
		fmt.Printf("err TextGet1: %v\n", err)
		return body, errors.New("网络请求失败")
	}

	if resp.StatusCode != 200 { //如果状态码错误 就打印错误 并且返回空
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return body, errors.New("状态码异常")
	}

	defer resp.Body.Close()                 //关闭连接
	body, err2 := ioutil.ReadAll(resp.Body) //将获取到的数据使用 ioutil.ReadAll 函数 转化一下
	if err2 != nil {
		fmt.Printf("err TextGet2: %v body %v\n", err2)
		return body, errors.New("数据转换为字符串错误 一般是返回的数据是压缩的 需要在请求头里面设置接口明码返回")
		//	设置接口明码返回  head["Accept-Encoding"] = "identity"
	}

	return body, nil

}

//get请求示例
func http_get() {
	//get请求示例
	head := make(map[string]string, 0)
	head["Accept-Encoding"] = "identity"
	head["user-agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"
	head["referer"] = "https://zh.stripchat.com/Baby-Nina"
	head["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	head["accept-language"] = "zh-CN,zh;q=0.9,ja;q=0.8"
	head["upgrade-insecure-requests"] = "1"

	// s, err := TextGet("https://www.tokyomotion.net/user/yutong/videos", head) //请求普通的html网页

	// s, err := ByteGet("https://cdn.tokyo-motion.net/img/webapp-icon.png", head) //下载文件

	s, _ := JsonGet("http://api.jktab.xyz/up_info/list/?Page=1&pageSize=155", head) //请求json数据
	test := s["data"].([]interface{})                                               //转换为可以遍历的数组
	for _, val := range test {
		m := val.(map[string]interface{}) //转换为可以提取的map
		fmt.Printf("val: %v\n", m)
		fmt.Printf("val: %v\n", m["name"])
	}

}

//post请求示例
func http_post() {
	// link := "http://httpbin.org/post"
	// link := "http://127.0.0.1:5000/post"

	head := make(map[string]string, 0)
	head["Accept-Encoding"] = "identity"
	head["user-agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"
	head["referer"] = "https://zh.stripchat.com/Baby-Nina"
	head["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	head["accept-language"] = "zh-CN,zh;q=0.9,ja;q=0.8"
	head["upgrade-insecure-requests"] = "1"

	data := make(map[string]string, 0)
	data["name"] = "海虎"
	data["age"] = "48"

	// s, err := TextPost(link, head, data)   //提交form表单 返回字符串类型
	// m, err := JsonPost(link, head, data) //提交form表单 返回json类型
	//f, err :=  BytePost(link, head, data)              //提交form表单 数据流 用于下载文件

	// s, err := Post_Json_Text(link, head, data)    //提交json数据 返回字符串类型
	// m, err := Post_Json_Json(link, head, data) //提交json数据 返回json数据
	// Post_Json_Byte(link, head, data) //提交json数据 返回数据流 用于下载文件

}

func main() {
	// http_get()
	http_post()
}

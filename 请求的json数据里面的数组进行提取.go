package main

import (
	"fmt"
)

func main() {
	// body:=
	s := make(map[string]interface{})

	//这里的s模拟请求返回的数据
	//现在我需要提取json里面的数组数据 进行遍历 然后再提取内容

	test := s["data"].([]interface{}) //将里面的data字段提取出来  然后使用.([]interface{})  转换为可以for循环的数组

	for _, val := range test {
		m := val.(map[string]interface{})  //再将遍历出来的内容 转换为map数组 这样才可以用 m["字段名"]  来进行取值
		fmt.Printf("val: %v\n", m["name"]) //这样就行了 用m[""] 就可以提取指定值了

	}

}

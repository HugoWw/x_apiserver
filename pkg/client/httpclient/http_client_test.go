package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestHttpGet(t *testing.T) {
	//u, err := url.Parse("http://www.baidu.com:9090/api?index=1")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Printf("u.Path:%s ,u.RawQuery:%s ,u.Host:%s, u.Scheme:%s, u.RawFragment:%s ,u.Fragment:%s ,u.OmitHost:%v, u.Opaque:%s  \n",
	//	u.Path, u.RawQuery, u.Host, u.Scheme, u.RawFragment, u.Fragment, u.OmitHost, u.Opaque)

	type testData1 struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Bodys  string `json:"body"`
	}

	client, err := NewHttpClient(http.DefaultClient, 30, "https://jsonplaceholder.typicode.com", 5)
	if err != nil {
		fmt.Println("new http client error:", err)
		return
	}

	data1 := testData1{}
	err = client.Get().Prefix("/v1").SetPath("/posts/1").Do(context.TODO()).Into(&data1)
	if err != nil {
		fmt.Println("http client request err:", err)
		return
	}

	fmt.Printf("out put get data info:%+v\n", data1)

}

func TestHttpPost(t *testing.T) {
	type testData1 struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Bodys  string `json:"body"`
	}

	data1 := testData1{
		UserId: 10086,
		Id:     101,
		Title:  "PostTest",
		Bodys:  "hello Post test",
	}

	client, err := NewHttpClient(http.DefaultClient, 30, "https://jsonplaceholder.typicode.com", 5)
	if err != nil {
		fmt.Println("new http client error:", err)
		return
	}

	b, _ := json.Marshal(data1)

	data2 := testData1{}
	err = client.Post().Body(b).SetHeader("Content-type", "application/json; charset=UTF-8").SetPath("/posts").Do(context.TODO()).Into(&data2)
	if err != nil {
		fmt.Println("http client post err:", err)
		return
	}

	fmt.Printf("out put get data info:%+v\n", data2)
}

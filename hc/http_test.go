package hc

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	params := map[string]interface{}{
		"1": "2",
	}
	resp, _ := Get("http://httpbin.org/get", params)

	fmt.Println(resp.AsString())

	m := make(map[string]interface{})

	resp.AsJson(&m)

	fmt.Println(fmt.Sprint(m))

}

func TestPostForm(t *testing.T) {
	params := map[string]interface{}{
		"1": "2",
	}
	resp, _ := PostForm("http://httpbin.org/post", params)

	fmt.Println(resp.AsString())

	m := make(map[string]interface{})

	resp.AsJson(&m)

	fmt.Println(fmt.Sprint(m))
}

func TestPostJson(t *testing.T) {
	params := map[string]interface{}{
		"1": "2",
	}
	resp, _ := PostJson("http://httpbin.org/post", params)

	fmt.Println(resp.AsString())

	m := make(map[string]interface{})

	resp.AsJson(&m)

	fmt.Println(fmt.Sprint(m))
}
